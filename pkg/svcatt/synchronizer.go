package svcatt

import (
	"fmt"

	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
	coreinformers "k8s.io/client-go/informers/core/v1"
	coreclient "k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templateinformers "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/templates/experimental"
	templateslisters "github.com/Azure/service-catalog-templates/pkg/client/listers/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcatinformers "github.com/kubernetes-incubator/service-catalog/pkg/client/informers_generated/externalversions/servicecatalog/v1beta1"
	svcatlisters "github.com/kubernetes-incubator/service-catalog/pkg/client/listers_generated/servicecatalog/v1beta1"
)

const (
	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by the Templates controller"
)

type Synchronizer struct {
	resolver *resolver

	coreClient      coreclient.Interface
	templatesClient templatesclient.Interface
	svcatClient     svcatclient.Interface

	secretLister          corelisters.SecretLister
	instanceLister        templateslisters.TemplatedInstanceLister
	bindingLister         templateslisters.TemplatedBindingLister
	bindingTemplateLister templateslisters.BindingTemplateLister
	svcatInstanceLister   svcatlisters.ServiceInstanceLister
	svcatBindingLister    svcatlisters.ServiceBindingLister
}

func NewSynchronizer(coreClient coreclient.Interface, templatesClient templatesclient.Interface, svcatClient svcatclient.Interface,
	coreInformers coreinformers.Interface, templatesInformers templateinformers.Interface, svcatInformers svcatinformers.Interface) *Synchronizer {
	return &Synchronizer{
		coreClient:            coreClient,
		templatesClient:       templatesClient,
		svcatClient:           svcatClient,
		secretLister:          coreInformers.Secrets().Lister(),
		instanceLister:        templatesInformers.TemplatedInstances().Lister(),
		bindingLister:         templatesInformers.TemplatedBindings().Lister(),
		bindingTemplateLister: templatesInformers.BindingTemplates().Lister(),
		svcatInstanceLister:   svcatInformers.ServiceInstances().Lister(),
		svcatBindingLister:    svcatInformers.ServiceBindings().Lister(),
		resolver:              newResolver(templatesClient, svcatClient),
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
		_, err := s.bindingLister.TemplatedBindings(object.GetNamespace()).Get(owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of %s '%s'", object.GetSelfLink(), owner.Kind, owner.Name)
			return false
		}
		return true
	case templates.InstanceKind:
		_, err := s.instanceLister.TemplatedInstances(object.GetNamespace()).Get(owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of %s '%s'", object.GetSelfLink(), owner.Kind, owner.Name)
			return false
		}
		return true
	case "ServiceBinding":
		// Lookup the binding that owns the resource
		svcBnd, err := s.svcatBindingLister.ServiceBindings(object.GetNamespace()).Get(owner.Name)
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

	// Get the TemplatedInstance resource with this namespace/name
	cachedInst, err := s.instanceLister.TemplatedInstances(namespace).Get(name)
	if err != nil {
		// The TemplatedInstance resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("instance '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}
	// TODO: Figure out the best practices for proactive DeepCopies and avoiding pointers so I don't mess up the cache
	inst := cachedInst.DeepCopy()

	instanceName := inst.Name
	if instanceName == "" {
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
	svcInst, err := s.svcatInstanceLister.ServiceInstances(inst.Namespace).Get(instanceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		cachedTemplate, err := s.resolver.ResolveInstanceTemplate(*inst)
		if err != nil {
			// TODO: Update status to unresolvable
			return false, inst, err
		}
		// TODO: Figure out the best practices for proactive DeepCopies and avoiding pointers so I don't mess up the cache
		template := cachedTemplate.DeepCopy()

		// Apply changes from the template to the instance
		inst, err = ApplyInstanceTemplate(*inst, *template)
		if err != nil {
			return false, inst, err
		}
		inst, err = s.templatesClient.TemplatesExperimental().TemplatedInstances(inst.Namespace).Update(inst)
		if err != nil {
			return false, inst, err
		}

		// Convert the templated resource into a service catalog resource
		svcInst, err = BuildServiceInstance(*inst)
		if err != nil {
			return false, inst, err
		}
		svcInst, err = s.svcatClient.ServicecatalogV1beta1().ServiceInstances(inst.Namespace).Create(svcInst)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, inst, err
	}

	// If the service instance is not controlled by our shadow instance, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(svcInst, inst) {
		msg := fmt.Sprintf(MessageResourceExists, svcInst.Name)
		return false, inst, fmt.Errorf(msg)
	}

	// TODO: Detect when the plan must be re-resolved

	// If this number of the replicas on the TemplatedInstance resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if inst.Spec.Parameters != nil && (svcInst.Spec.Parameters == nil || string(inst.Spec.Parameters.Raw) != string(svcInst.Spec.Parameters.Raw)) {
		glog.V(4).Infof("Syncing instance %s back to service instance %s", inst.SelfLink, svcInst.SelfLink)
		svcInst = RefreshServiceInstance(inst, svcInst)
		svcInst, err = s.svcatClient.ServicecatalogV1beta1().ServiceInstances(svcInst.Namespace).Update(svcInst)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, inst, err
	}

	//
	// Update shadow instance status with the service instance state
	//
	// Finally, we update the status block of the TemplatedInstance resource to reflect the
	// current state of the world
	err = s.updateInstanceStatus(inst, svcInst)
	if err != nil {
		return false, inst, err
	}

	return true, inst, nil
}

func (s *Synchronizer) updateInstanceStatus(inst *templates.TemplatedInstance, svcInst *svcat.ServiceInstance) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	instCopy := inst.DeepCopy()
	// TODO: add resolved fields to the status
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the TemplatedInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templatesClient.TemplatesExperimental().TemplatedInstances(inst.Namespace).Update(instCopy)
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

	// Get the resource with this namespace/name
	cachedBnd, err := s.bindingLister.TemplatedBindings(namespace).Get(name)
	if err != nil {
		// The resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("binding '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}
	bnd := cachedBnd.DeepCopy()

	bindingName := bnd.Name
	if bindingName == "" {
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
	svcBnd, err := s.svcatBindingLister.ServiceBindings(bnd.Namespace).Get(bindingName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		cachedTemplate, err := s.resolver.ResolveBindingTemplate(*bnd)
		if err != nil {
			// TODO: Update status to unresolvable
			return false, bnd, err
		}
		template := cachedTemplate.DeepCopy()

		// Apply changes from the template to the instance
		bnd, err = ApplyBindingTemplate(*bnd, *template)
		if err != nil {
			return false, bnd, err
		}
		bnd, err = s.templatesClient.TemplatesExperimental().TemplatedBindings(bnd.Namespace).Update(bnd)
		if err != nil {
			return false, bnd, err
		}

		// Convert the templated resource into a service catalog resource
		svcBnd = BuildServiceBinding(*bnd)
		svcBnd, err = s.svcatClient.ServicecatalogV1beta1().ServiceBindings(svcBnd.Namespace).Create(svcBnd)
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, bnd, err
	}

	// If the service catalog resource is not controlled by our shadow resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(svcBnd, bnd) {
		msg := fmt.Sprintf(MessageResourceExists, svcBnd.Name)
		return false, bnd, fmt.Errorf(msg)
	}

	//
	// Sync updates to shadow resource back to the service catalog resource
	//
	// TODO: sync other fields
	if bnd.Spec.Parameters != nil && (svcBnd.Spec.Parameters == nil || string(bnd.Spec.Parameters.Raw) != string(svcBnd.Spec.Parameters.Raw)) {
		glog.V(4).Infof("Syncing shadow binding %s back to service catalog binding %s", bnd.SelfLink, svcBnd.SelfLink)
		svcBnd = RefreshServiceBinding(bnd, svcBnd)
		svcBnd, err = s.svcatClient.ServicecatalogV1beta1().ServiceBindings(svcBnd.Namespace).Update(svcBnd)
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return false, bnd, err
	}

	//
	// Update shadow resource status with the service catalog resource state
	//
	err = s.updateBindingStatus(bnd, svcBnd)
	if err != nil {
		return false, bnd, err
	}

	return true, bnd, nil
}

func (s *Synchronizer) updateBindingStatus(bnd *templates.TemplatedBinding, svcBnd *svcat.ServiceBinding) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	bndCopy := bnd.DeepCopy()
	// TODO: add resolved fields to the status
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the TemplatedInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templatesClient.TemplatesExperimental().TemplatedBindings(bnd.Namespace).Update(bndCopy)
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

	// Get the resource with this namespace/name
	cachedSecret, err := s.secretLister.Secrets(namespace).Get(name)
	if err != nil {
		// The resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			util.HandleError(fmt.Errorf("secret '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}

	svcSecret := cachedSecret.DeepCopy()

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
	shadowSecretName := toShadowSecretName(svcSecret.Name)
	secret, err := s.secretLister.Secrets(svcSecret.Namespace).Get(shadowSecretName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		ownerSvcBnd := metav1.GetControllerOf(svcSecret)
		if ownerSvcBnd == nil {
			// Ignore unmanaged secrets
			return false, nil, nil
		}
		svcBnd, err := s.svcatBindingLister.ServiceBindings(svcSecret.Namespace).Get(ownerSvcBnd.Name)
		if err != nil {
			return false, svcSecret, err
		}

		ownerBnd := metav1.GetControllerOf(svcBnd)
		if ownerSvcBnd == nil {
			// Ignore unmanaged resources
			return false, nil, nil
		}
		bnd, err := s.bindingLister.TemplatedBindings(svcBnd.Namespace).Get(ownerBnd.Name)
		if err != nil {
			return false, svcSecret, err
		}

		secret, err = BuildShadowSecret(svcSecret, *bnd)
		if err != nil {
			return false, svcSecret, err
		}
		secret, err = s.coreClient.CoreV1().Secrets(secret.Namespace).Create(secret)
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
	if refreshedSecret, changed := RefreshSecret(*svcSecret, *secret); changed {
		secret, err = s.coreClient.CoreV1().Secrets(refreshedSecret.Namespace).Update(refreshedSecret)

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

func (s *Synchronizer) updateSecretStatus(secret *core.Secret, svcSecret *core.Secret) error {
	// TODO: do I need to update the binding instead of the secret to note successful synchronization?
	return nil
}
