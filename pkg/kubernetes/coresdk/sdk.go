package coresdk

import (
	"fmt"

	"github.com/golang/glog"
	corefactory "k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	coreclient "k8s.io/client-go/kubernetes"
	coreinterfaces "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// SDK wrapper around the generated Go client for the Service Catalog Templates API
type SDK struct {
	Client  coreclient.Interface
	Factory corefactory.SharedInformerFactory

	informers    coreinformers.Interface
	secretLister corelisters.SecretLister
}

func New(client coreclient.Interface, factory corefactory.SharedInformerFactory) *SDK {
	return &SDK{
		Client:  client,
		Factory: factory,
	}
}

func (sdk *SDK) Init(stopCh <-chan struct{}) error {
	go sdk.Factory.Start(stopCh)
	secretsInformer := sdk.Cache().Secrets().Informer()
	if ok := cache.WaitForCacheSync(stopCh,
		secretsInformer.HasSynced); !ok {
		return fmt.Errorf("failed to wait for core caches to sync")
	}
	glog.Info("Finished synchronizing core caches")
	return nil
}

// Core is the underlying generated Core versioned interface
// It should be used instead of accessing the client directly.
func (sdk *SDK) Core() coreinterfaces.CoreV1Interface {
	return sdk.Client.CoreV1()
}

func (sdk *SDK) Cache() coreinformers.Interface {
	if sdk.informers == nil {
		sdk.informers = sdk.Factory.Core().V1()
	}
	return sdk.informers
}

func (sdk *SDK) SecretCache() corelisters.SecretLister {
	if sdk.secretLister == nil {
		sdk.secretLister = sdk.Cache().Secrets().Lister()
	}
	return sdk.secretLister
}
