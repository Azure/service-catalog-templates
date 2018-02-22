package svcatt

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
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

	templatesClient templatesclient.Interface
	svcatClient     svcatclient.Interface

	instanceLister      templateslisters.CatalogInstanceLister
	bindingLister       templateslisters.CatalogBindingLister
	svcatInstanceLister svcatlisters.ServiceInstanceLister
	svcatBindingLister  svcatlisters.ServiceBindingLister
}

func NewSynchronizer(templatesClient templatesclient.Interface, svcatClient svcatclient.Interface,
	templatesInformers templateinformers.Interface, svcatInformers svcatinformers.Interface) *Synchronizer {
	return &Synchronizer{
		templatesClient:     templatesClient,
		svcatClient:         svcatClient,
		instanceLister:      templatesInformers.CatalogInstances().Lister(),
		bindingLister:       templatesInformers.CatalogBindings().Lister(),
		svcatInstanceLister: svcatInformers.ServiceInstances().Lister(),
		svcatBindingLister:  svcatInformers.ServiceBindings().Lister(),
		resolver:            newResolver(templatesClient, svcatClient),
	}
}

// IsManaged determines if a resource is managed by a shadow resource.
func (s *Synchronizer) IsManaged(object metav1.Object) bool {
	owner := metav1.GetControllerOf(object)
	if owner == nil {
		// Ignore unmanaged service catalog bindings
		return false
	}

	// Try to retrieve the binding that is shadowing the service catalog binding
	switch owner.Kind {
	case templates.BindingKind:
		_, err := s.bindingLister.CatalogBindings(object.GetNamespace()).Get(owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of binding '%s'", object.GetSelfLink(), owner.Name)
			return false
		}
		return true
	case templates.InstanceKind:
		_, err := s.instanceLister.CatalogInstances(object.GetNamespace()).Get(owner.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of instance '%s'", object.GetSelfLink(), owner.Name)
			return false
		}
		return true
	}

	return false
}

// SynchronizeInstance accepts an instance key (namespace/name)
// and attempts to synchronize it with a service catalog instance.
// * ok - Synchronization was successful.
// * instance - The instance resource.
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

	// Get the CatalogInstance resource with this namespace/name
	cachedInst, err := s.instanceLister.CatalogInstances(namespace).Get(name)
	if err != nil {
		// The CatalogInstance resource may no longer exist, in which case we stop
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

		svcInst, err = BuildServiceInstance(*inst, *template)
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

	// If this number of the replicas on the CatalogInstance resource is specified, and the
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
	// Finally, we update the status block of the CatalogInstance resource to reflect the
	// current state of the world
	err = s.updateInstanceStatus(inst, svcInst)
	if err != nil {
		return false, inst, err
	}

	return true, inst, nil
}

func (s *Synchronizer) updateInstanceStatus(inst *templates.CatalogInstance, svcInst *svcat.ServiceInstance) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	instCopy := inst.DeepCopy()
	//fooCopy.Status.Message = deployment.Status.Message
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the CatalogInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templatesClient.TemplatesExperimental().CatalogInstances(inst.Namespace).Update(instCopy)
	return err
}

// SynchronizeInstance accepts an instance key (namespace/name)
// and attempts to synchronize it with a service catalog instance.
// * ok - Synchronization was successful.
// * instance - The instance resource.
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
	cachedBnd, err := s.bindingLister.CatalogBindings(namespace).Get(name)
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

		svcBnd, err = BuildServiceBinding(*bnd, *template)
		if err != nil {
			return false, bnd, err
		}
		svcBnd, err = s.svcatClient.ServicecatalogV1beta1().ServiceBindings(bnd.Namespace).Create(svcBnd)
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

func (s *Synchronizer) updateBindingStatus(bnd *templates.CatalogBinding, svcBnd *svcat.ServiceBinding) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	bndCopy := bnd.DeepCopy()
	//fooCopy.Status.Message = deployment.Status.Message
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the CatalogInstance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templatesClient.TemplatesExperimental().CatalogBindings(bnd.Namespace).Update(bndCopy)
	return err
}
