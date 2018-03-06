package sdk

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templatesinterfaces "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/typed/templates/experimental"
	templatesfactory "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions"
	templatesinformer "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/templates/experimental"
	templateslisters "github.com/Azure/service-catalog-templates/pkg/client/listers/templates/experimental"
)

// SDK wrapper around the generated Go client for the Service Catalog Templates API
type SDK struct {
	Client  templatesclient.Interface
	Factory templatesfactory.SharedInformerFactory

	informers               templatesinformer.Interface
	templatedInstanceLister templateslisters.TemplatedInstanceLister
	templatedBindingLister  templateslisters.TemplatedBindingLister
}

func New(client templatesclient.Interface, factory templatesfactory.SharedInformerFactory) *SDK {
	return &SDK{
		Client:  client,
		Factory: factory,
	}
}

func (sdk *SDK) Init(stopCh <-chan struct{}) error {

	tbnd := sdk.Cache().TemplatedBindings().Informer()
	tinst := sdk.Cache().TemplatedInstances().Informer()

	go sdk.Factory.Start(stopCh)

	if ok := cache.WaitForCacheSync(stopCh,
		tbnd.HasSynced,
		tinst.HasSynced); !ok {
		return fmt.Errorf("failed to wait for templates caches to sync")
	}

	glog.Info("Finished synchronizing templates caches")
	return nil
}

// Templates is the underlying generated Templates versioned interface
// It should be used instead of accessing the client directly.
func (sdk *SDK) Templates() templatesinterfaces.TemplatesExperimentalInterface {
	return sdk.Client.TemplatesExperimental()
}

func (sdk *SDK) Cache() templatesinformer.Interface {
	if sdk.informers == nil {
		sdk.informers = sdk.Factory.Templates().Experimental()
	}
	return sdk.informers
}

func (sdk *SDK) InstanceCache() templateslisters.TemplatedInstanceLister {
	if sdk.templatedInstanceLister == nil {
		sdk.templatedInstanceLister = sdk.Cache().TemplatedInstances().Lister()
	}
	return sdk.templatedInstanceLister
}

func (sdk *SDK) BindingCache() templateslisters.TemplatedBindingLister {
	if sdk.templatedBindingLister == nil {
		sdk.templatedBindingLister = sdk.Cache().TemplatedBindings().Lister()
	}
	return sdk.templatedBindingLister
}

func (sdk *SDK) filterByServiceTypeLabel(serviceType string) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		templates.FieldServiceTypeName: serviceType,
	})
}
