// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	scheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// InstanceTemplatesGetter has a method to return a InstanceTemplateInterface.
// A group's client should implement this interface.
type InstanceTemplatesGetter interface {
	InstanceTemplates(namespace string) InstanceTemplateInterface
}

// InstanceTemplateInterface has methods to work with InstanceTemplate resources.
type InstanceTemplateInterface interface {
	Create(*experimental.InstanceTemplate) (*experimental.InstanceTemplate, error)
	Update(*experimental.InstanceTemplate) (*experimental.InstanceTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.InstanceTemplate, error)
	List(opts v1.ListOptions) (*experimental.InstanceTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.InstanceTemplate, err error)
	InstanceTemplateExpansion
}

// instanceTemplates implements InstanceTemplateInterface
type instanceTemplates struct {
	client rest.Interface
	ns     string
}

// newInstanceTemplates returns a InstanceTemplates
func newInstanceTemplates(c *TemplatesExperimentalClient, namespace string) *instanceTemplates {
	return &instanceTemplates{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the instanceTemplate, and returns the corresponding instanceTemplate object, and an error if there is any.
func (c *instanceTemplates) Get(name string, options v1.GetOptions) (result *experimental.InstanceTemplate, err error) {
	result = &experimental.InstanceTemplate{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instancetemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of InstanceTemplates that match those selectors.
func (c *instanceTemplates) List(opts v1.ListOptions) (result *experimental.InstanceTemplateList, err error) {
	result = &experimental.InstanceTemplateList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("instancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested instanceTemplates.
func (c *instanceTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("instancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a instanceTemplate and creates it.  Returns the server's representation of the instanceTemplate, and an error, if there is any.
func (c *instanceTemplates) Create(instanceTemplate *experimental.InstanceTemplate) (result *experimental.InstanceTemplate, err error) {
	result = &experimental.InstanceTemplate{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("instancetemplates").
		Body(instanceTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a instanceTemplate and updates it. Returns the server's representation of the instanceTemplate, and an error, if there is any.
func (c *instanceTemplates) Update(instanceTemplate *experimental.InstanceTemplate) (result *experimental.InstanceTemplate, err error) {
	result = &experimental.InstanceTemplate{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("instancetemplates").
		Name(instanceTemplate.Name).
		Body(instanceTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the instanceTemplate and deletes it. Returns an error if one occurs.
func (c *instanceTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("instancetemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *instanceTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("instancetemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched instanceTemplate.
func (c *instanceTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.InstanceTemplate, err error) {
	result = &experimental.InstanceTemplate{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("instancetemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
