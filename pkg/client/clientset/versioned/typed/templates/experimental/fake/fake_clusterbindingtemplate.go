// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
package fake

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterBindingTemplates implements ClusterBindingTemplateInterface
type FakeClusterBindingTemplates struct {
	Fake *FakeTemplatesExperimental
}

var clusterbindingtemplatesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "clusterbindingtemplates"}

var clusterbindingtemplatesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "ClusterBindingTemplate"}

// Get takes name of the clusterBindingTemplate, and returns the corresponding clusterBindingTemplate object, and an error if there is any.
func (c *FakeClusterBindingTemplates) Get(name string, options v1.GetOptions) (result *experimental.ClusterBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterbindingtemplatesResource, name), &experimental.ClusterBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterBindingTemplate), err
}

// List takes label and field selectors, and returns the list of ClusterBindingTemplates that match those selectors.
func (c *FakeClusterBindingTemplates) List(opts v1.ListOptions) (result *experimental.ClusterBindingTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterbindingtemplatesResource, clusterbindingtemplatesKind, opts), &experimental.ClusterBindingTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.ClusterBindingTemplateList{}
	for _, item := range obj.(*experimental.ClusterBindingTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterBindingTemplates.
func (c *FakeClusterBindingTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterbindingtemplatesResource, opts))
}

// Create takes the representation of a clusterBindingTemplate and creates it.  Returns the server's representation of the clusterBindingTemplate, and an error, if there is any.
func (c *FakeClusterBindingTemplates) Create(clusterBindingTemplate *experimental.ClusterBindingTemplate) (result *experimental.ClusterBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterbindingtemplatesResource, clusterBindingTemplate), &experimental.ClusterBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterBindingTemplate), err
}

// Update takes the representation of a clusterBindingTemplate and updates it. Returns the server's representation of the clusterBindingTemplate, and an error, if there is any.
func (c *FakeClusterBindingTemplates) Update(clusterBindingTemplate *experimental.ClusterBindingTemplate) (result *experimental.ClusterBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterbindingtemplatesResource, clusterBindingTemplate), &experimental.ClusterBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterBindingTemplate), err
}

// Delete takes name of the clusterBindingTemplate and deletes it. Returns an error if one occurs.
func (c *FakeClusterBindingTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterbindingtemplatesResource, name), &experimental.ClusterBindingTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterBindingTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterbindingtemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.ClusterBindingTemplateList{})
	return err
}

// Patch applies the patch and returns the patched clusterBindingTemplate.
func (c *FakeClusterBindingTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterbindingtemplatesResource, name, data, subresources...), &experimental.ClusterBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterBindingTemplate), err
}
