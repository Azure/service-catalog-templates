// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Azure/service-catalog-templates/pkg/kubernetes/core-sdk"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-sdk"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates-sdk"
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
	configure()

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
	svcatSDK := servicecatalogsdk.New(svcatClient, svcatInformerFactory)
	templateSDK := servicecatalogtempltesdk.New(templatesClient, templatesInformerFactory, svcatSDK)

	// Wait for the caches to be synced before starting
	glog.Info("Initializing...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	initG, ctx := errgroup.WithContext(ctx)
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

func configure() {
	flag.Parse()

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			localConfig := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
			if _, err := os.Stat(localConfig); err == nil {
				kubeconfig = localConfig
			}
		}
	}

	if kubeconfig == "" {
		glog.Infof("Using kubeconfig: %s", kubeconfig)
	} else {
		glog.Info("Using in-cluster kubeconfig")
	}
}
