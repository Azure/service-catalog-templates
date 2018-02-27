package sdk

import (
	"fmt"

	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templatesinterfaces "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/typed/templates/experimental"
	templatesfactory "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions"
	templatesinformer "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/templates/experimental"
	templateslisters "github.com/Azure/service-catalog-templates/pkg/client/listers/templates/experimental"
	"k8s.io/client-go/tools/cache"
)

// SDK wrapper around the generated Go client for the Service Catalog Templates API
type SDK struct {
	Client  templatesclient.Interface
	Factory templatesfactory.SharedInformerFactory

	informers               templatesinformer.Interface
	bindingTemplateLister   templateslisters.BindingTemplateLister
	instanceTemplateLister  templateslisters.InstanceTemplateLister
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
	go sdk.Factory.Start(stopCh)
	if ok := cache.WaitForCacheSync(stopCh,
		// TODO: These should probably be saved and reused
		sdk.Cache().InstanceTemplates().Informer().HasSynced,
		sdk.Cache().BindingTemplates().Informer().HasSynced,
		sdk.Cache().TemplatedBindings().Informer().HasSynced,
		sdk.Cache().TemplatedInstances().Informer().HasSynced); !ok {
		return fmt.Errorf("failed to wait for informer caches to sync")
	}
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

func (sdk *SDK) InstanceTemplateCache() templateslisters.InstanceTemplateLister {
	if sdk.instanceTemplateLister == nil {
		sdk.instanceTemplateLister = sdk.Cache().InstanceTemplates().Lister()
	}
	return sdk.instanceTemplateLister
}

func (sdk *SDK) BindingCache() templateslisters.TemplatedBindingLister {
	if sdk.templatedBindingLister == nil {
		sdk.templatedBindingLister = sdk.Cache().TemplatedBindings().Lister()
	}
	return sdk.templatedBindingLister
}

func (sdk *SDK) BindingTemplateCache() templateslisters.BindingTemplateLister {
	if sdk.bindingTemplateLister == nil {
		sdk.bindingTemplateLister = sdk.Cache().BindingTemplates().Lister()
	}
	return sdk.bindingTemplateLister
}
