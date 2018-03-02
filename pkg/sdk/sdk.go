package sdk

import (
	"fmt"

	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	templatesinterfaces "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/typed/templates/experimental"
	templatesfactory "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions"
	templatesinformer "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/templates/experimental"
	templateslisters "github.com/Azure/service-catalog-templates/pkg/client/listers/templates/experimental"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/cache"
)

// SDK wrapper around the generated Go client for the Service Catalog Templates API
type SDK struct {
	Client  templatesclient.Interface
	Factory templatesfactory.SharedInformerFactory

	informers                     templatesinformer.Interface
	bindingTemplateLister         templateslisters.BindingTemplateLister
	clusterBindingTemplateLister  templateslisters.ClusterBindingTemplateLister
	brokerBindingTemplateLister   templateslisters.BrokerBindingTemplateLister
	instanceTemplateLister        templateslisters.InstanceTemplateLister
	clusterInstanceTemplateLister templateslisters.ClusterInstanceTemplateLister
	brokerInstanceTemplateLister  templateslisters.BrokerInstanceTemplateLister
	templatedInstanceLister       templateslisters.TemplatedInstanceLister
	templatedBindingLister        templateslisters.TemplatedBindingLister
}

func New(client templatesclient.Interface, factory templatesfactory.SharedInformerFactory) *SDK {
	return &SDK{
		Client:  client,
		Factory: factory,
	}
}

func (sdk *SDK) Init(stopCh <-chan struct{}) error {
	go sdk.Factory.Start(stopCh)
	instt := sdk.Cache().InstanceTemplates().Informer()
	cinstt := sdk.Cache().ClusterInstanceTemplates().Informer()
	binstt := sdk.Cache().BrokerInstanceTemplates().Informer()
	bndt := sdk.Cache().BindingTemplates().Informer()
	cbndt := sdk.Cache().ClusterBindingTemplates().Informer()
	bbndt := sdk.Cache().BrokerBindingTemplates().Informer()
	tbnd := sdk.Cache().TemplatedBindings().Informer()
	tinst := sdk.Cache().TemplatedInstances().Informer()
	if ok := cache.WaitForCacheSync(stopCh,
		// TODO: These should probably be saved and reused
		instt.HasSynced,
		cinstt.HasSynced,
		binstt.HasSynced,
		bndt.HasSynced,
		cbndt.HasSynced,
		bbndt.HasSynced,
		tbnd.HasSynced,
		tinst.HasSynced); !ok {
		return fmt.Errorf("failed to wait for templates caches to sync")
	}

	glog.Info("Finished synchronizing template caches")
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

func (sdk *SDK) ClusterInstanceTemplateCache() templateslisters.ClusterInstanceTemplateLister {
	if sdk.clusterInstanceTemplateLister == nil {
		sdk.clusterInstanceTemplateLister = sdk.Cache().ClusterInstanceTemplates().Lister()
	}
	return sdk.clusterInstanceTemplateLister
}

func (sdk *SDK) BrokerInstanceTemplateCache() templateslisters.BrokerInstanceTemplateLister {
	if sdk.brokerInstanceTemplateLister == nil {
		sdk.brokerInstanceTemplateLister = sdk.Cache().BrokerInstanceTemplates().Lister()
	}
	return sdk.brokerInstanceTemplateLister
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

func (sdk *SDK) ClusterBindingTemplateCache() templateslisters.ClusterBindingTemplateLister {
	if sdk.clusterBindingTemplateLister == nil {
		sdk.clusterBindingTemplateLister = sdk.Cache().ClusterBindingTemplates().Lister()
	}
	return sdk.clusterBindingTemplateLister
}

func (sdk *SDK) BrokerBindingTemplateCache() templateslisters.BrokerBindingTemplateLister {
	if sdk.brokerBindingTemplateLister == nil {
		sdk.brokerBindingTemplateLister = sdk.Cache().BrokerBindingTemplates().Lister()
	}
	return sdk.brokerBindingTemplateLister
}
