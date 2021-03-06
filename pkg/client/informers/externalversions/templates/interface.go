// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

// This file was automatically generated by informer-gen

package templates

import (
	internalinterfaces "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/internalinterfaces"
	experimental "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/templates/experimental"
)

// Interface provides access to each of this group's versions.
type Interface interface {
	// Experimental provides access to shared informers for resources in Experimental.
	Experimental() experimental.Interface
}

type group struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Experimental returns a new experimental.Interface.
func (g *group) Experimental() experimental.Interface {
	return experimental.New(g.factory, g.namespace, g.tweakListOptions)
}
