// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtempltesdk

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/pkg/service-catalog-sdk"
	"github.com/golang/glog"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	svcatSDK                *servicecatalogsdk.SDK
	informers               templatesinformer.Interface
	templatedInstanceLister templateslisters.TemplatedInstanceLister
	templatedBindingLister  templateslisters.TemplatedBindingLister
}

func New(client templatesclient.Interface, factory templatesfactory.SharedInformerFactory, svcatSDK *servicecatalogsdk.SDK) *SDK {
	return &SDK{
		Client:   client,
		Factory:  factory,
		svcatSDK: svcatSDK,
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

func filterByServiceTypeLabel(serviceType string) meta.ListOptions {
	opts := meta.ListOptions{}
	if serviceType != "" {
		opts.LabelSelector = labels.SelectorFromSet(map[string]string{
			templates.FieldServiceTypeName: serviceType,
		}).String()
	}
	return opts
}
