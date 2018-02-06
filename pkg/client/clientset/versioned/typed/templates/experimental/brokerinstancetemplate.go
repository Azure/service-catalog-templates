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

// BrokerInstanceTemplatesGetter has a method to return a BrokerInstanceTemplateInterface.
// A group's client should implement this interface.
type BrokerInstanceTemplatesGetter interface {
	BrokerInstanceTemplates() BrokerInstanceTemplateInterface
}

// BrokerInstanceTemplateInterface has methods to work with BrokerInstanceTemplate resources.
type BrokerInstanceTemplateInterface interface {
	Create(*experimental.BrokerInstanceTemplate) (*experimental.BrokerInstanceTemplate, error)
	Update(*experimental.BrokerInstanceTemplate) (*experimental.BrokerInstanceTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.BrokerInstanceTemplate, error)
	List(opts v1.ListOptions) (*experimental.BrokerInstanceTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BrokerInstanceTemplate, err error)
	BrokerInstanceTemplateExpansion
}

// brokerInstanceTemplates implements BrokerInstanceTemplateInterface
type brokerInstanceTemplates struct {
	client rest.Interface
}

// newBrokerInstanceTemplates returns a BrokerInstanceTemplates
func newBrokerInstanceTemplates(c *TemplatesExperimentalClient) *brokerInstanceTemplates {
	return &brokerInstanceTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the brokerInstanceTemplate, and returns the corresponding brokerInstanceTemplate object, and an error if there is any.
func (c *brokerInstanceTemplates) Get(name string, options v1.GetOptions) (result *experimental.BrokerInstanceTemplate, err error) {
	result = &experimental.BrokerInstanceTemplate{}
	err = c.client.Get().
		Resource("brokerinstancetemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BrokerInstanceTemplates that match those selectors.
func (c *brokerInstanceTemplates) List(opts v1.ListOptions) (result *experimental.BrokerInstanceTemplateList, err error) {
	result = &experimental.BrokerInstanceTemplateList{}
	err = c.client.Get().
		Resource("brokerinstancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested brokerInstanceTemplates.
func (c *brokerInstanceTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("brokerinstancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a brokerInstanceTemplate and creates it.  Returns the server's representation of the brokerInstanceTemplate, and an error, if there is any.
func (c *brokerInstanceTemplates) Create(brokerInstanceTemplate *experimental.BrokerInstanceTemplate) (result *experimental.BrokerInstanceTemplate, err error) {
	result = &experimental.BrokerInstanceTemplate{}
	err = c.client.Post().
		Resource("brokerinstancetemplates").
		Body(brokerInstanceTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a brokerInstanceTemplate and updates it. Returns the server's representation of the brokerInstanceTemplate, and an error, if there is any.
func (c *brokerInstanceTemplates) Update(brokerInstanceTemplate *experimental.BrokerInstanceTemplate) (result *experimental.BrokerInstanceTemplate, err error) {
	result = &experimental.BrokerInstanceTemplate{}
	err = c.client.Put().
		Resource("brokerinstancetemplates").
		Name(brokerInstanceTemplate.Name).
		Body(brokerInstanceTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the brokerInstanceTemplate and deletes it. Returns an error if one occurs.
func (c *brokerInstanceTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("brokerinstancetemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *brokerInstanceTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("brokerinstancetemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched brokerInstanceTemplate.
func (c *brokerInstanceTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BrokerInstanceTemplate, err error) {
	result = &experimental.BrokerInstanceTemplate{}
	err = c.client.Patch(pt).
		Resource("brokerinstancetemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
