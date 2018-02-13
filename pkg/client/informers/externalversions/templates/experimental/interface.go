// This file was automatically generated by informer-gen

package experimental

import (
	internalinterfaces "github.com/Azure/service-catalog-templates/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Instances returns a InstanceInformer.
	Instances() InstanceInformer
	// InstanceTemplates returns a InstanceTemplateInformer.
	InstanceTemplates() InstanceTemplateInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Instances returns a InstanceInformer.
func (v *version) Instances() InstanceInformer {
	return &instanceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// InstanceTemplates returns a InstanceTemplateInformer.
func (v *version) InstanceTemplates() InstanceTemplateInformer {
	return &instanceTemplateInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
