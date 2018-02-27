package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	util "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers"
	coreclient "k8s.io/client-go/kubernetes"
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
	// SuccessSynced is used as part of the Event 'reason' when a shadow resource is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a shadow resource fails
	// to sync due to an unmanaged resource of the same name already existing.
	ErrResourceExists = "ErrResourceExists"
	// MessageResourceSynced is the message used for an Event fired when a shadow resource
	// is synced successfully
	MessageResourceSynced = "Shadow resource synced successfully"
)

// Controller is the controller implementation for Instance resources
// NOTE: This is the stock CRD implementation from https://github.com/kubernetes/sample-controller
// all interesting logic should live in ../svcatt/synchronizer.go
type Controller struct {
	synchronizer *svcatt.Synchronizer

	coreClient           coreclient.Interface
	instancesSynced      cache.InformerSynced
	svcatInstancesSynced cache.InformerSynced
	bindingsSynced       cache.InformerSynced
	svcatBindingsSynced  cache.InformerSynced
	secretsSynced        cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	instanceQ workqueue.RateLimitingInterface
	bindingQ  workqueue.RateLimitingInterface
	secretQ   workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	coreClient coreclient.Interface,
	svcatClient svcatclient.Interface,
	templatesClient templatesclient.Interface,
	coreInformerFactory coreinformers.SharedInformerFactory,
	svcatInformerFactory svcatinformers.SharedInformerFactory,
	templatesInformerFactory templatesinformers.SharedInformerFactory) *Controller {

	// obtain references to shared index templatesinformers for the Deployment and Instance
	// types.
	templatesInformers := templatesInformerFactory.Templates().Experimental()
	instanceInformer := templatesInformers.TemplatedInstances().Informer()
	bindingInformer := templatesInformers.TemplatedBindings().Informer()
	svcatInformers := svcatInformerFactory.Servicecatalog().V1beta1()
	svcatInstanceInformer := svcatInformers.ServiceInstances().Informer()
	svcatBindingInformer := svcatInformers.ServiceBindings().Informer()
	coreInformers := coreInformerFactory.Core().V1()
	secretInformer := coreInformers.Secrets().Informer()

	// Create event broadcaster
	// Add service-catalog-templates-controller types to the default Kubernetes Scheme so Events can be
	// logged for service-catalog-templates-controller types.
	templatesscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: coreClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	c := &Controller{
		coreClient:           coreClient,
		synchronizer:         svcatt.NewSynchronizer(coreClient, templatesClient, svcatClient, coreInformers, templatesInformers, svcatInformers),
		instancesSynced:      instanceInformer.HasSynced,
		bindingsSynced:       bindingInformer.HasSynced,
		secretsSynced:        secretInformer.HasSynced,
		svcatInstancesSynced: svcatInstanceInformer.HasSynced,
		svcatBindingsSynced:  svcatBindingInformer.HasSynced,
		instanceQ:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Instances"),
		bindingQ:             workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Bindings"),
		secretQ:              workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Secrets"),
		recorder:             recorder,
	}

	glog.Info("Setting up event handlers")

	// Set up an event handler for when shadow resources change
	instanceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueueResource(obj, c.instanceQ)
		},
		UpdateFunc: func(old, new interface{}) {
			c.enqueueResource(new, c.instanceQ)
		},
	})
	bindingInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueueResource(obj, c.bindingQ)
		},
		UpdateFunc: func(old, new interface{}) {
			c.enqueueResource(new, c.bindingQ)
		},
	})
	secretInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueueResource(obj, c.secretQ)
		},
		UpdateFunc: func(old, new interface{}) {
			c.enqueueResource(new, c.secretQ)
		},
	})

	// Set up an event handler for when managed resources change. This
	// handler will lookup the owner of the given resource, and if it is
	// owned by a shadow resource will enqueue that resource for
	// processing. This way, we don't need to implement custom logic for
	// handling managed resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	svcatInstanceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.handleManagedResource,
		UpdateFunc: func(old, new interface{}) {
			newInst := new.(*svcat.ServiceInstance)
			oldInst := old.(*svcat.ServiceInstance)
			if newInst.ResourceVersion == oldInst.ResourceVersion {
				// Periodic resync will send update events for all known instances.
				// Two different versions of the same instance will always have different RVs.
				return
			}
			c.handleManagedResource(new)
		},
		DeleteFunc: c.handleManagedResource,
	})
	svcatBindingInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.handleManagedResource,
		UpdateFunc: func(old, new interface{}) {
			newBnd := new.(*svcat.ServiceBinding)
			oldBnd := old.(*svcat.ServiceBinding)
			if newBnd.ResourceVersion == oldBnd.ResourceVersion {
				// Periodic resync will send update events for all known instances.
				// Two different versions of the same instance will always have different RVs.
				return
			}
			c.handleManagedResource(new)
		},
		DeleteFunc: c.handleManagedResource,
	})

	return c
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer util.HandleCrash()
	defer c.instanceQ.ShutDown()
	defer c.bindingQ.ShutDown()
	defer c.secretQ.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Templates controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh,
		c.secretsSynced,
		c.svcatInstancesSynced, c.instancesSynced,
		c.svcatBindingsSynced, c.bindingsSynced); !ok {
		return fmt.Errorf("failed to wait for informer caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(func() {
			for c.processNextWorkItem(c.instanceQ, c.synchronizer.SynchronizeInstance) {
			}
		}, time.Second, stopCh)
		go wait.Until(func() {
			for c.processNextWorkItem(c.bindingQ, c.synchronizer.SynchronizeBinding) {
			}
		}, time.Second, stopCh)
		go wait.Until(func() {
			for c.processNextWorkItem(c.secretQ, c.synchronizer.SynchronizeSecret) {
			}
		}, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

type resourceSynchronizationHandler func(key string) (bool, runtime.Object, error)

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem(q workqueue.RateLimitingInterface, sync resourceSynchronizationHandler) bool {
	obj, shutdown := q.Get()

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
		defer q.Done(obj)
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
			c.instanceQ.Forget(obj)
			util.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the sync, passing it the namespace/name string of the
		// Instance resource to be synced.
		if err := c.synchronizeResource(key, sync); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		q.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		util.HandleError(err)
		return true
	}

	return true
}

// synchronizeResource compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the resource
// with the current status of the resource.
func (c *Controller) synchronizeResource(key string, sync resourceSynchronizationHandler) error {
	ok, obj, err := sync(key)
	if err != nil {
		// Append a warning to the resource
		if obj != nil {
			c.recorder.Event(obj, corev1.EventTypeWarning, ErrResourceExists, err.Error())
		}
		return err
	}

	// Record a successful sync
	//
	if ok {
		glog.Infof("Successfully synced %s '%s'", obj.GetObjectKind().GroupVersionKind().Kind, key)
		c.recorder.Event(obj, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	}

	return nil
}

// enqueueResource takes a resource and converts it into a namespace/name
// string which is then put onto the specified work queue.
func (c *Controller) enqueueResource(obj interface{}, q workqueue.RateLimitingInterface) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		util.HandleError(err)
		return
	}
	q.AddRateLimited(key)
}

// handleManagedResource will take any resource implementing metav1.Object and attempt
// to find the shadow resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleManagedResource(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			util.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			util.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ok := c.synchronizer.IsManaged(object); ok {
		c.enqueueResource(obj, c.instanceQ)
		return
	}
}
