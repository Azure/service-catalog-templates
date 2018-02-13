package svcatt

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templateslisters "github.com/Azure/service-catalog-templates/pkg/client/listers/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
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
}

func NewSynchronizer(templatesClient templatesclient.Interface, svcatClient svcatclient.Interface,
	instanceLister templateslisters.InstanceLister, svcatInstancesLister svcatlisters.ServiceInstanceLister) *Synchronizer {
	return &Synchronizer{
		templatesClient:      templatesClient,
		svcatClient:          svcatClient,
		instanceLister:       instanceLister,
		svcatInstancesLister: svcatInstancesLister,
	}
}

// SynchronizeInstance accepts an instance key (namespace/name)
// and attempts to synchronize it with a service catalog instance.
// * ok - Synchronization was successful.
// * instance - The instance resource.
// * error - Fatal synchronization error.
func (s *Synchronizer) SynchronizeInstance(key string) (retry bool, instance *templates.Instance, err error) {
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
	inst, err := s.instanceLister.Instances(namespace).Get(name)
	if err != nil {
		// The Instance resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("instance '%s' in work queue no longer exists", key))
			return false, nil, nil
		}

		return false, nil, err
	}

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
		svcInst = BuildServiceInstance(inst, nil)
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

	/*
		// If this number of the replicas on the Instance resource is specified, and the
		// number does not equal the current desired replicas on the Deployment, we
		// should update the Deployment resource.
		if foo.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
			glog.V(4).Infof("Instance %s replicas: %d, deployment replicas: %d", name, *foo.Spec.Replicas, *deployment.Spec.Replicas)
			deployment, err = c.kubeClient.AppsV1().Deployments(foo.Namespace).Update(newInstance(foo))
		}

		// If an error occurs during Update, we'll requeue the item so we can
		// attempt processing again later. THis could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}
	*/

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
