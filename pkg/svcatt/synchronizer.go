package svcatt

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/pkg/kubernetes/coresdk"
	"github.com/Azure/service-catalog-templates/pkg/sdk"
	"github.com/Azure/service-catalog-templates/pkg/svcatsdk"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/Azure/service-catalog-templates/pkg/svcatt/builder"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

const (
	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by the Templates controller"
)

type Synchronizer struct {
	resolver    *resolver
	coreSDK     *coresdk.SDK
	templateSDK *sdk.SDK
	svcatSDK    *svcatsdk.SDK
}

func NewSynchronizer(coreSDK *coresdk.SDK, templateSDK *sdk.SDK, svcatSDK *svcatsdk.SDK) *Synchronizer {
	return &Synchronizer{
		coreSDK:     coreSDK,
		templateSDK: templateSDK,
		svcatSDK:    svcatSDK,
		resolver:    newResolver(templateSDK, svcatSDK),
	}
}

// IsManaged determines if a resource is managed by a shadow resource.
func (s *Synchronizer) IsManaged(object metav1.Object) bool {
	owner := metav1.GetControllerOf(object)
	if owner == nil {
		// Ignore unmanaged service catalog resources
		return false
	}

	// Try to retrieve the resource that is shadowing the service catalog resource
	switch owner.Kind {
	case templates.BindingKind:
		_, err := s.templateSDK.GetBindingFromCache(object.GetNamespace(), owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of %s '%s'", object.GetSelfLink(), owner.Kind, owner.Name)
			return false
		}
		return true
	case templates.InstanceKind:
		_, err := s.templateSDK.GetInstanceFromCache(object.GetNamespace(), owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of %s '%s'", object.GetSelfLink(), owner.Kind, owner.Name)
			return false
		}
		return true
	case "ServiceBinding":
		// Lookup the binding that owns the resource
		svcBnd, err := s.svcatSDK.GetBindingFromCache(object.GetNamespace(), owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of %s '%s'", object.GetSelfLink(), owner.Kind, owner.Name)
			return false
		}

		// The binding must be owned by the templates controller
		return s.IsManaged(svcBnd)
	}

	return false
}

// SynchronizeInstance accepts an instance key (namespace/name)
// and attempts to synchronize it with a service catalog instance.
// * ok - Synchronization was successful.
// * resource - The resource.
// * error - Fatal synchronization error.
func (s *Synchronizer) SynchronizeInstance(key string) (bool, runtime.Object, error) {
	//
	// Get shadow instance
	//

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		util.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return false, nil, nil
	}

	tinst, err := s.templateSDK.GetInstanceFromCache(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("instance '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}

	if tinst.Name == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		util.HandleError(fmt.Errorf("%s: instance name must be specified", key))
		return false, nil, nil
	}

	//
	// Sync shadow to service catalog instance
	//

	// Get the corresponding service instance from the service catalog
	svcInst, err := s.svcatSDK.GetInstanceFromCache(tinst.Namespace, tinst.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		template, err := s.resolver.ResolveInstanceTemplate(*tinst)
		if err != nil {
			// TODO: Update status to unresolvable
			return false, tinst, err
		}

		// Apply changes from the template to the instance
		tinst, err = builder.ApplyInstanceTemplate(*tinst, *template)
		if err != nil {
			return false, tinst, err
		}
		tinst, err = s.templateSDK.Templates().TemplatedInstances(tinst.Namespace).Update(tinst)
		if err != nil {
			return false, tinst, err
		}

		// Convert the templated resource into a service catalog resource
		svcInst, err = builder.BuildServiceInstance(*tinst)
		if err != nil {
			return false, tinst, err
		}
		svcInst, err = s.svcatSDK.ServiceCatalog().ServiceInstances(tinst.Namespace).Create(svcInst)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, tinst, err
	}

	// If the service instance is not controlled by our shadow instance, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(svcInst, tinst) {
		msg := fmt.Sprintf(MessageResourceExists, svcInst.Name)
		return false, tinst, fmt.Errorf(msg)
	}

	// TODO: Detect when the plan must be re-resolved

	// If this number of the replicas on the TemplatedInstance resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if tinst.Spec.Parameters != nil && (svcInst.Spec.Parameters == nil || string(tinst.Spec.Parameters.Raw) != string(svcInst.Spec.Parameters.Raw)) {
		glog.V(4).Infof("Syncing instance %s back to service instance %s", tinst.SelfLink, svcInst.SelfLink)
		svcInst = builder.RefreshServiceInstance(tinst, svcInst)
		svcInst, err = s.svcatSDK.ServiceCatalog().ServiceInstances(svcInst.Namespace).Update(svcInst)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, tinst, err
	}

	//
	// Update shadow instance status with the service instance state
	//
	// Finally, we update the status block of the TemplatedInstance resource to reflect the
	// current state of the world
	err = s.updateInstanceStatus(tinst, svcInst)
	if err != nil {
		return false, tinst, err
	}

	return true, tinst, nil
}

func (s *Synchronizer) updateInstanceStatus(inst *templates.TemplatedInstance, svcInst *svcat.ServiceInstance) error {
	// TODO: add resolved fields to the status
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the TemplatedInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templateSDK.Templates().TemplatedInstances(inst.Namespace).Update(inst)
	return err
}

// SynchronizeBinding accepts an binding key (namespace/name)
// and attempts to synchronize it with a service catalog binding.
// * ok - Synchronization was successful.
// * resource - The resource.
// * error - Fatal synchronization error.
func (s *Synchronizer) SynchronizeBinding(key string) (bool, runtime.Object, error) {
	//
	// Get shadow resource
	//

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		util.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return false, nil, nil
	}

	tbnd, err := s.templateSDK.GetBindingFromCache(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("binding '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}

	if tbnd.Name == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		util.HandleError(fmt.Errorf("%s: binding name must be specified", key))
		return false, nil, nil
	}

	//
	// Sync shadow resource back to service catalog resource
	//

	// Get the corresponding service catalog resource
	svcBnd, err := s.svcatSDK.GetBindingFromCache(tbnd.Namespace, tbnd.Name)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		template, err := s.resolver.ResolveBindingTemplate(*tbnd)
		if err != nil {
			// TODO: Update status to unresolvable
			return false, tbnd, err
		}

		// Apply changes from the template to the instance
		tbnd, err = builder.ApplyBindingTemplate(*tbnd, *template)
		if err != nil {
			return false, tbnd, err
		}
		tbnd, err = s.templateSDK.Templates().TemplatedBindings(tbnd.Namespace).Update(tbnd)
		if err != nil {
			return false, tbnd, err
		}

		// Convert the templated resource into a service catalog resource
		svcBnd = builder.BuildServiceBinding(*tbnd)
		svcBnd, err = s.svcatSDK.ServiceCatalog().ServiceBindings(svcBnd.Namespace).Create(svcBnd)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, tbnd, err
	}

	// If the service catalog resource is not controlled by our shadow resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(svcBnd, tbnd) {
		msg := fmt.Sprintf(MessageResourceExists, svcBnd.Name)
		return false, tbnd, fmt.Errorf(msg)
	}

	//
	// Sync updates to shadow resource back to the service catalog resource
	//
	// TODO: sync other fields
	if tbnd.Spec.Parameters != nil && (svcBnd.Spec.Parameters == nil || string(tbnd.Spec.Parameters.Raw) != string(svcBnd.Spec.Parameters.Raw)) {
		glog.V(4).Infof("Syncing shadow binding %s back to service catalog binding %s", tbnd.SelfLink, svcBnd.SelfLink)
		svcBnd = builder.RefreshServiceBinding(tbnd, svcBnd)
		svcBnd, err = s.svcatSDK.ServiceCatalog().ServiceBindings(svcBnd.Namespace).Update(svcBnd)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, tbnd, err
	}

	//
	// Update shadow resource status with the service catalog resource state
	//
	err = s.updateBindingStatus(tbnd, svcBnd)
	if err != nil {
		return false, tbnd, err
	}

	return true, tbnd, nil
}

func (s *Synchronizer) updateBindingStatus(bnd *templates.TemplatedBinding, svcBnd *svcat.ServiceBinding) error {
	// TODO: add resolved fields to the status
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the TemplatedInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templateSDK.Templates().TemplatedBindings(bnd.Namespace).Update(bnd)
	return err
}

// SynchronizeSecret accepts a secret key (namespace/name)
// and attempts to synchronize the bound secret with the template secret.
// * ok - Synchronization was successful.
// * resource - The resource.
// * error - Fatal synchronization error.
func (s *Synchronizer) SynchronizeSecret(key string) (bool, runtime.Object, error) {
	//
	// Get shadow resource
	//
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		util.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return false, nil, nil
	}

	svcSecret, err := s.coreSDK.GetSecretFromCache(namespace, name)
	if err != nil {
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("secret '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}

	if svcSecret.Name == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		util.HandleError(fmt.Errorf("%s: secret name must be specified", key))
		return false, nil, nil
	}

	//
	// Sync service catalog resource back to the shadow resource
	//

	// Get the corresponding shadow resource
	shadowSecretName := builder.BoundSecretName(svcSecret.Name)
	secret, err := s.coreSDK.GetSecretFromCache(svcSecret.Namespace, shadowSecretName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		tbnd, err := s.GetTemplatedBindingFromShadowSecret(svcSecret)
		if err != nil {
			return false, svcSecret, err
		}
		if tbnd == nil {
			// ignore unmanaged secrets
			return false, nil, nil
		}

		secret, err = builder.BuildBoundSecret(*svcSecret, *tbnd)
		if err != nil {
			return false, svcSecret, err
		}
		secret, err = s.coreSDK.Core().Secrets(secret.Namespace).Create(secret)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, svcSecret, err
	}

	// If the shadow secret is not controlled by the service catalog managed secret,
	// we should log a warning to the event recorder and retry
	if !metav1.IsControlledBy(secret, svcSecret) {
		return false, nil, nil
	}

	//
	// Sync updates to service catalog resource back to the shadow resource
	//
	tbnd, err := s.GetTemplatedBindingFromShadowSecret(svcSecret)
	if err != nil {
		return false, svcSecret, err
	}
	if tbnd == nil {
		// ignore unmanaged secrets
		return false, nil, nil
	}

	if refreshedSecret, changed := builder.RefreshSecret(*svcSecret, *tbnd, *secret); changed {
		secret, err = s.coreSDK.Core().Secrets(refreshedSecret.Namespace).Update(refreshedSecret)

		// If an error occurs during Update, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return false, svcSecret, err
		}
	}

	//
	// Update shadow resource status with the service catalog resource state
	//
	err = s.updateSecretStatus(secret, svcSecret)
	if err != nil {
		return false, svcSecret, err
	}

	return true, svcSecret, nil
}

func (s *Synchronizer) GetTemplatedBindingFromShadowSecret(svcSecret *core.Secret) (*templates.TemplatedBinding, error) {
	svcBnd, err := s.svcatSDK.GetSecretOwner(svcSecret)
	if err != nil {
		return nil, err
	}

	if svcBnd == nil {
		return nil, nil
	}

	return s.templateSDK.GetBindingOwner(svcBnd)
}

func (s *Synchronizer) updateSecretStatus(secret *core.Secret, svcSecret *core.Secret) error {
	// TODO: do I need to update the binding instead of the secret to note successful synchronization?
	return nil
}
