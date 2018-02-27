package main

import (
	"flag"
	"time"

	"github.com/Azure/service-catalog-templates/pkg/kubernetes/coresdk"
	"github.com/Azure/service-catalog-templates/pkg/sdk"
	"github.com/Azure/service-catalog-templates/pkg/svcatsdk"
	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	coreinformers "k8s.io/client-go/informers"
	coreclient "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"

	clientset "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	informers "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions"
	"github.com/Azure/service-catalog-templates/pkg/controller"
	"github.com/Azure/service-catalog-templates/pkg/signals"
	svcatclientset "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcatinformers "github.com/kubernetes-incubator/service-catalog/pkg/client/informers_generated/externalversions"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	coreClient, err := coreclient.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	svcatClient, err := svcatclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building service catalog clientset: %s", err.Error())
	}

	templatesClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building example clientset: %s", err.Error())
	}

	duration := time.Second * 30
	coreInformerFactory := coreinformers.NewSharedInformerFactory(coreClient, duration)
	svcatInformerFactory := svcatinformers.NewSharedInformerFactory(svcatClient, duration)
	templatesInformerFactory := informers.NewSharedInformerFactory(templatesClient, duration)

	coreSDK := coresdk.New(coreClient, coreInformerFactory)
	templateSDK := sdk.New(templatesClient, templatesInformerFactory)
	svcatSDK := svcatsdk.New(svcatClient, svcatInformerFactory)

	// Wait for the caches to be synced before starting
	glog.Info("Initializing...")
	var initG errgroup.Group
	initG.Go(func() error { return coreSDK.Init(stopCh) })
	initG.Go(func() error { return svcatSDK.Init(stopCh) })
	initG.Go(func() error { return templateSDK.Init(stopCh) })
	if err := initG.Wait(); err != nil {
		glog.Fatalf("Error initializing informer caches: %s", err)
	}

	controller := controller.NewController(coreSDK, templateSDK, svcatSDK)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
