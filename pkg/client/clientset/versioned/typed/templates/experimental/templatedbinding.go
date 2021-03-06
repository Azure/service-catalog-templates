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

// TemplatedBindingsGetter has a method to return a TemplatedBindingInterface.
// A group's client should implement this interface.
type TemplatedBindingsGetter interface {
	TemplatedBindings(namespace string) TemplatedBindingInterface
}

// TemplatedBindingInterface has methods to work with TemplatedBinding resources.
type TemplatedBindingInterface interface {
	Create(*experimental.TemplatedBinding) (*experimental.TemplatedBinding, error)
	Update(*experimental.TemplatedBinding) (*experimental.TemplatedBinding, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.TemplatedBinding, error)
	List(opts v1.ListOptions) (*experimental.TemplatedBindingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedBinding, err error)
	TemplatedBindingExpansion
}

// templatedBindings implements TemplatedBindingInterface
type templatedBindings struct {
	client rest.Interface
	ns     string
}

// newTemplatedBindings returns a TemplatedBindings
func newTemplatedBindings(c *TemplatesExperimentalClient, namespace string) *templatedBindings {
	return &templatedBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the templatedBinding, and returns the corresponding templatedBinding object, and an error if there is any.
func (c *templatedBindings) Get(name string, options v1.GetOptions) (result *experimental.TemplatedBinding, err error) {
	result = &experimental.TemplatedBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("templatedbindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TemplatedBindings that match those selectors.
func (c *templatedBindings) List(opts v1.ListOptions) (result *experimental.TemplatedBindingList, err error) {
	result = &experimental.TemplatedBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("templatedbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested templatedBindings.
func (c *templatedBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("templatedbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a templatedBinding and creates it.  Returns the server's representation of the templatedBinding, and an error, if there is any.
func (c *templatedBindings) Create(templatedBinding *experimental.TemplatedBinding) (result *experimental.TemplatedBinding, err error) {
	result = &experimental.TemplatedBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("templatedbindings").
		Body(templatedBinding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a templatedBinding and updates it. Returns the server's representation of the templatedBinding, and an error, if there is any.
func (c *templatedBindings) Update(templatedBinding *experimental.TemplatedBinding) (result *experimental.TemplatedBinding, err error) {
	result = &experimental.TemplatedBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("templatedbindings").
		Name(templatedBinding.Name).
		Body(templatedBinding).
		Do().
		Into(result)
	return
}

// Delete takes name of the templatedBinding and deletes it. Returns an error if one occurs.
func (c *templatedBindings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("templatedbindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *templatedBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("templatedbindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched templatedBinding.
func (c *templatedBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedBinding, err error) {
	result = &experimental.TemplatedBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("templatedbindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
