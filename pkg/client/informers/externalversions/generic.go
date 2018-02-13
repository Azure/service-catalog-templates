// This file was automatically generated by informer-gen

package externalversions

import (
	"fmt"

	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=templates.servicecatalog.k8s.io, Version=experimental
	case experimental.SchemeGroupVersion.WithResource("instances"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Templates().Experimental().Instances().Informer()}, nil
	case experimental.SchemeGroupVersion.WithResource("instancetemplates"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Templates().Experimental().InstanceTemplates().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
