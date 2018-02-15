package svcatt

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
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
	MessageResourceExists = "Resource %q already exists and is not managed by Instance"
)

type Synchronizer struct {
	templatesClient      templatesclient.Interface
	svcatClient          svcatclient.Interface
	instanceLister       templateslisters.InstanceLister
	svcatInstancesLister svcatlisters.ServiceInstanceLister
	resolver             *resolver
}

func NewSynchronizer(templatesClient templatesclient.Interface, svcatClient svcatclient.Interface,
	templatesInformers templateinformers.Interface, svcatInformers svcatinformers.Interface) *Synchronizer {
	return &Synchronizer{
		templatesClient:      templatesClient,
		svcatClient:          svcatClient,
		instanceLister:       templatesInformers.Instances().Lister(),
		svcatInstancesLister: svcatInformers.ServiceInstances().Lister(),
		resolver:             newResolver(templatesClient, svcatClient),
	}
}

// IsManagedInstance determines if a resource is managed by a shadow instance.
func (s *Synchronizer) IsManagedInstance(object metav1.Object) (bool, *templates.Instance) {
	owner := metav1.GetControllerOf(object)
	if owner == nil {
		return false, nil
	}

	// Ignore unmanaged service catalog instances
	if owner.Kind != templates.InstanceKind {
		return false, nil
	}

	// Try to retrieve the instance that is shadowing the service catalog instance
	instance, err := s.instanceLister.Instances(object.GetNamespace()).Get(owner.Name)
	if err != nil {
		glog.V(4).Infof("ignoring orphaned object '%s' of instance '%s'", object.GetSelfLink(), owner.Name)
		return false, nil
	}

	return true, instance
}

// SynchronizeInstance accepts an instance key (namespace/name)
// and attempts to synchronize it with a service catalog instance.
// * ok - Synchronization was successful.
// * instance - The instance resource.
// * error - Fatal synchronization error.
func (s *Synchronizer) SynchronizeInstance(key string) (bool, *templates.Instance, error) {
	//
	// Get shadow instance
	//

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return false, nil, nil
	}

	// Get the Instance resource with this namespace/name
	cachedInst, err := s.instanceLister.Instances(namespace).Get(name)
	if err != nil {
		// The Instance resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("instance '%s' in work queue no longer exists", key))
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
		runtime.HandleError(fmt.Errorf("%s: instance name must be specified", key))
		return false, nil, nil
	}

	//
	// Sync shadow to service catalog instance
	//

	// Get the corresponding service instance from the service catalog
	svcInst, err := s.svcatInstancesLister.ServiceInstances(inst.Namespace).Get(instanceName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		cachedTemplate, err := s.resolver.ResolveTemplate(*inst)
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

	// If this number of the replicas on the Instance resource is specified, and the
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
	// Finally, we update the status block of the Instance resource to reflect the
	// current state of the world
	err = s.updateInstanceStatus(inst, svcInst)
	if err != nil {
		return false, inst, err
	}

	return true, inst, nil
}

func (s *Synchronizer) updateInstanceStatus(inst *templates.Instance, svcInst *svcat.ServiceInstance) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	instCopy := inst.DeepCopy()
	//fooCopy.Status.Message = deployment.Status.Message
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the Instance resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := s.templatesClient.TemplatesExperimental().Instances(inst.Namespace).Update(instCopy)
	return err
}
