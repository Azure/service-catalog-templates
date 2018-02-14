package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templatesscheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	templatesinformers "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions"
	"github.com/Azure/service-catalog-templates/pkg/svcatt"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcatinformers "github.com/kubernetes-incubator/service-catalog/pkg/client/informers_generated/externalversions"
)

const controllerAgentName = "service-catalog-templates"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Instance is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Instance fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceSynced is the message used for an Event fired when a Instance
	// is synced successfully
	MessageResourceSynced = "Instance synced successfully"
)

// Controller is the controller implementation for Instance resources
// NOTE: This is the stock CRD implementation from https://github.com/kubernetes/sample-controller
// all interesting logic should live in ../svcatt/synchronizer.go
type Controller struct {
	synchronizer *svcatt.Synchronizer

	kubeClient           kubernetes.Interface
	instancesSynced      cache.InformerSynced
	svcatInstancesSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	kubeClient kubernetes.Interface,
	svcatClient svcatclient.Interface,
	templatesClient templatesclient.Interface,
	svcatInformerFactory svcatinformers.SharedInformerFactory,
	templatesInformerFactory templatesinformers.SharedInformerFactory) *Controller {

	// obtain references to shared index templatesinformers for the Deployment and Instance
	// types.
	instanceInformer := templatesInformerFactory.Templates().Experimental().Instances()
	svcatInstanceInformer := svcatInformerFactory.Servicecatalog().V1beta1().ServiceInstances()

	// Create event broadcaster
	// Add service-catalog-templates-controller types to the default Kubernetes Scheme so Events can be
	// logged for service-catalog-templates-controller types.
	templatesscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClient:           kubeClient,
		synchronizer:         svcatt.NewSynchronizer(templatesClient, svcatClient, instanceInformer.Lister(), svcatInstanceInformer.Lister()),
		instancesSynced:      instanceInformer.Informer().HasSynced,
		svcatInstancesSynced: svcatInstanceInformer.Informer().HasSynced,
		workqueue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Instances"),
		recorder:             recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Instance resources change
	instanceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueInstance,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueInstance(new)
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Instance resource will enqueue that Instance resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	svcatInstanceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newInst := new.(*svcat.ServiceInstance)
			oldInst := old.(*svcat.ServiceInstance)
			if newInst.ResourceVersion == oldInst.ResourceVersion {
				// Periodic resync will send update events for all known instances.
				// Two different versions of the same instance will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Templates controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.svcatInstancesSynced, c.instancesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process Instance resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Instance resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Instance resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	ok, instance, err := c.synchronizer.SynchronizeInstance(key)
	if err != nil {
		// Append a warning to the instance
		if instance != nil {
			c.recorder.Event(instance, corev1.EventTypeWarning, ErrResourceExists, err.Error())
		}
		return err
	}

	// Record a successful sync
	//
	if ok {
		c.recorder.Event(instance, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	}

	return nil
}

// enqueueInstance takes a Instance resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Instance.
func (c *Controller) enqueueInstance(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Instance resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Instance resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ok, instance := c.synchronizer.IsManagedInstance(object); ok {
		c.enqueueInstance(instance)
		return
	}
}
