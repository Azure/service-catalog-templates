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

// ClusterBindingTemplatesGetter has a method to return a ClusterBindingTemplateInterface.
// A group's client should implement this interface.
type ClusterBindingTemplatesGetter interface {
	ClusterBindingTemplates() ClusterBindingTemplateInterface
}

// ClusterBindingTemplateInterface has methods to work with ClusterBindingTemplate resources.
type ClusterBindingTemplateInterface interface {
	Create(*experimental.ClusterBindingTemplate) (*experimental.ClusterBindingTemplate, error)
	Update(*experimental.ClusterBindingTemplate) (*experimental.ClusterBindingTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.ClusterBindingTemplate, error)
	List(opts v1.ListOptions) (*experimental.ClusterBindingTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterBindingTemplate, err error)
	ClusterBindingTemplateExpansion
}

// clusterBindingTemplates implements ClusterBindingTemplateInterface
type clusterBindingTemplates struct {
	client rest.Interface
}

// newClusterBindingTemplates returns a ClusterBindingTemplates
func newClusterBindingTemplates(c *TemplatesExperimentalClient) *clusterBindingTemplates {
	return &clusterBindingTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterBindingTemplate, and returns the corresponding clusterBindingTemplate object, and an error if there is any.
func (c *clusterBindingTemplates) Get(name string, options v1.GetOptions) (result *experimental.ClusterBindingTemplate, err error) {
	result = &experimental.ClusterBindingTemplate{}
	err = c.client.Get().
		Resource("clusterbindingtemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterBindingTemplates that match those selectors.
func (c *clusterBindingTemplates) List(opts v1.ListOptions) (result *experimental.ClusterBindingTemplateList, err error) {
	result = &experimental.ClusterBindingTemplateList{}
	err = c.client.Get().
		Resource("clusterbindingtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterBindingTemplates.
func (c *clusterBindingTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("clusterbindingtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a clusterBindingTemplate and creates it.  Returns the server's representation of the clusterBindingTemplate, and an error, if there is any.
func (c *clusterBindingTemplates) Create(clusterBindingTemplate *experimental.ClusterBindingTemplate) (result *experimental.ClusterBindingTemplate, err error) {
	result = &experimental.ClusterBindingTemplate{}
	err = c.client.Post().
		Resource("clusterbindingtemplates").
		Body(clusterBindingTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterBindingTemplate and updates it. Returns the server's representation of the clusterBindingTemplate, and an error, if there is any.
func (c *clusterBindingTemplates) Update(clusterBindingTemplate *experimental.ClusterBindingTemplate) (result *experimental.ClusterBindingTemplate, err error) {
	result = &experimental.ClusterBindingTemplate{}
	err = c.client.Put().
		Resource("clusterbindingtemplates").
		Name(clusterBindingTemplate.Name).
		Body(clusterBindingTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterBindingTemplate and deletes it. Returns an error if one occurs.
func (c *clusterBindingTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterbindingtemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterBindingTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("clusterbindingtemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterBindingTemplate.
func (c *clusterBindingTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterBindingTemplate, err error) {
	result = &experimental.ClusterBindingTemplate{}
	err = c.client.Patch(pt).
		Resource("clusterbindingtemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
