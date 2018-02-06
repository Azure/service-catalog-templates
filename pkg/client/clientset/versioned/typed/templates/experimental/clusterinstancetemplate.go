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

// ClusterInstanceTemplatesGetter has a method to return a ClusterInstanceTemplateInterface.
// A group's client should implement this interface.
type ClusterInstanceTemplatesGetter interface {
	ClusterInstanceTemplates() ClusterInstanceTemplateInterface
}

// ClusterInstanceTemplateInterface has methods to work with ClusterInstanceTemplate resources.
type ClusterInstanceTemplateInterface interface {
	Create(*experimental.ClusterInstanceTemplate) (*experimental.ClusterInstanceTemplate, error)
	Update(*experimental.ClusterInstanceTemplate) (*experimental.ClusterInstanceTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.ClusterInstanceTemplate, error)
	List(opts v1.ListOptions) (*experimental.ClusterInstanceTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterInstanceTemplate, err error)
	ClusterInstanceTemplateExpansion
}

// clusterInstanceTemplates implements ClusterInstanceTemplateInterface
type clusterInstanceTemplates struct {
	client rest.Interface
}

// newClusterInstanceTemplates returns a ClusterInstanceTemplates
func newClusterInstanceTemplates(c *TemplatesExperimentalClient) *clusterInstanceTemplates {
	return &clusterInstanceTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterInstanceTemplate, and returns the corresponding clusterInstanceTemplate object, and an error if there is any.
func (c *clusterInstanceTemplates) Get(name string, options v1.GetOptions) (result *experimental.ClusterInstanceTemplate, err error) {
	result = &experimental.ClusterInstanceTemplate{}
	err = c.client.Get().
		Resource("clusterinstancetemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterInstanceTemplates that match those selectors.
func (c *clusterInstanceTemplates) List(opts v1.ListOptions) (result *experimental.ClusterInstanceTemplateList, err error) {
	result = &experimental.ClusterInstanceTemplateList{}
	err = c.client.Get().
		Resource("clusterinstancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterInstanceTemplates.
func (c *clusterInstanceTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("clusterinstancetemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a clusterInstanceTemplate and creates it.  Returns the server's representation of the clusterInstanceTemplate, and an error, if there is any.
func (c *clusterInstanceTemplates) Create(clusterInstanceTemplate *experimental.ClusterInstanceTemplate) (result *experimental.ClusterInstanceTemplate, err error) {
	result = &experimental.ClusterInstanceTemplate{}
	err = c.client.Post().
		Resource("clusterinstancetemplates").
		Body(clusterInstanceTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a clusterInstanceTemplate and updates it. Returns the server's representation of the clusterInstanceTemplate, and an error, if there is any.
func (c *clusterInstanceTemplates) Update(clusterInstanceTemplate *experimental.ClusterInstanceTemplate) (result *experimental.ClusterInstanceTemplate, err error) {
	result = &experimental.ClusterInstanceTemplate{}
	err = c.client.Put().
		Resource("clusterinstancetemplates").
		Name(clusterInstanceTemplate.Name).
		Body(clusterInstanceTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the clusterInstanceTemplate and deletes it. Returns an error if one occurs.
func (c *clusterInstanceTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterinstancetemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterInstanceTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("clusterinstancetemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched clusterInstanceTemplate.
func (c *clusterInstanceTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterInstanceTemplate, err error) {
	result = &experimental.ClusterInstanceTemplate{}
	err = c.client.Patch(pt).
		Resource("clusterinstancetemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
