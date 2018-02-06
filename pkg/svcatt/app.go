// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcatt

import (
	"fmt"

	templatesclientset "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-sdk"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates-sdk"
	"github.com/kubernetes-incubator/service-catalog/pkg/svcat"
	"github.com/kubernetes-incubator/service-catalog/pkg/svcat/kube"
)

type ServiceCatalogApp = svcat.App

// App is the underlying application behind the svcat cli.
type App struct {
	*ServiceCatalogApp
	*servicecatalogtempltesdk.SDK

	// CurrentNamespace is the namespace set in the current context.
	CurrentNamespace string
}

// NewApp creates an svcat application.
func NewApp(kubeConfig, kubeContext string) (*App, error) {
	// Initialize a service catalog templates client
	cl, ns, err := getTemplatesClient(kubeConfig, kubeContext)
	if err != nil {
		return nil, err
	}

	svcApp, err := svcat.NewApp(kubeConfig, kubeContext)
	if err != nil {
		return nil, err
	}

	svcSDK := servicecatalogsdk.New(svcApp.ServiceCatalogClient, nil)
	app := &App{
		ServiceCatalogApp: svcApp,
		SDK:               servicecatalogtempltesdk.New(cl, nil, svcSDK),
		CurrentNamespace:  ns,
	}

	return app, nil
}

// getTemplatesClient creates a Service Catalog Templates config and client for a given kubeconfig context.
func getTemplatesClient(kubeConfig, kubeContext string) (client *templatesclientset.Clientset, namespaces string, err error) {
	config := kube.GetConfig(kubeContext, kubeConfig)

	currentNamespace, _, err := config.Namespace()
	if err != nil {
		return nil, "", fmt.Errorf("could not determine the namespace for the current context %q: %s", kubeContext, err)
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf("could not get Kubernetes config for context %q: %s", kubeContext, err)
	}

	client, err = templatesclientset.NewForConfig(restConfig)
	return client, currentNamespace, err
}
