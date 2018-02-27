package svcatsdk

import (
	"fmt"

	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcatfactory "github.com/kubernetes-incubator/service-catalog/pkg/client/informers_generated/externalversions"
	svcatinformers "github.com/kubernetes-incubator/service-catalog/pkg/client/informers_generated/externalversions/servicecatalog/v1beta1"
	svcatlisters "github.com/kubernetes-incubator/service-catalog/pkg/client/listers_generated/servicecatalog/v1beta1"
	svcatsdk "github.com/kubernetes-incubator/service-catalog/pkg/svcat/service-catalog"
	"k8s.io/client-go/tools/cache"
)

type SDK struct {
	*svcatsdk.SDK
	Factory svcatfactory.SharedInformerFactory

	informers      svcatinformers.Interface
	instanceLister svcatlisters.ServiceInstanceLister
	bindingLister  svcatlisters.ServiceBindingLister
}

func New(client svcatclient.Interface, factory svcatfactory.SharedInformerFactory) *SDK {
	return &SDK{
		SDK:     &svcatsdk.SDK{ServiceCatalogClient: client},
		Factory: factory,
	}
}

func (sdk *SDK) Init(stopCh <-chan struct{}) error {
	go sdk.Factory.Start(stopCh)
	if ok := cache.WaitForCacheSync(stopCh,
		// TODO: These should probably be saved and reused
		sdk.Cache().ServiceInstances().Informer().HasSynced,
		sdk.Cache().ServiceBindings().Informer().HasSynced); !ok {
		return fmt.Errorf("failed to wait for service catalog informer caches to sync")
	}
	return nil
}

func (sdk *SDK) Cache() svcatinformers.Interface {
	if sdk.informers == nil {
		sdk.informers = sdk.Factory.Servicecatalog().V1beta1()
	}
	return sdk.informers
}

func (sdk *SDK) InstanceCache() svcatlisters.ServiceInstanceLister {
	if sdk.instanceLister == nil {
		sdk.instanceLister = sdk.Cache().ServiceInstances().Lister()
	}
	return sdk.instanceLister
}

func (sdk *SDK) BindingCache() svcatlisters.ServiceBindingLister {
	if sdk.bindingLister == nil {
		sdk.bindingLister = sdk.Cache().ServiceBindings().Lister()
	}
	return sdk.bindingLister
}
