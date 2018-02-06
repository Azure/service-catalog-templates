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

// FakeBrokerBindingTemplates implements BrokerBindingTemplateInterface
type FakeBrokerBindingTemplates struct {
	Fake *FakeTemplatesExperimental
}

var brokerbindingtemplatesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "brokerbindingtemplates"}

var brokerbindingtemplatesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "BrokerBindingTemplate"}

// Get takes name of the brokerBindingTemplate, and returns the corresponding brokerBindingTemplate object, and an error if there is any.
func (c *FakeBrokerBindingTemplates) Get(name string, options v1.GetOptions) (result *experimental.BrokerBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(brokerbindingtemplatesResource, name), &experimental.BrokerBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerBindingTemplate), err
}

// List takes label and field selectors, and returns the list of BrokerBindingTemplates that match those selectors.
func (c *FakeBrokerBindingTemplates) List(opts v1.ListOptions) (result *experimental.BrokerBindingTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(brokerbindingtemplatesResource, brokerbindingtemplatesKind, opts), &experimental.BrokerBindingTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.BrokerBindingTemplateList{}
	for _, item := range obj.(*experimental.BrokerBindingTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested brokerBindingTemplates.
func (c *FakeBrokerBindingTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(brokerbindingtemplatesResource, opts))
}

// Create takes the representation of a brokerBindingTemplate and creates it.  Returns the server's representation of the brokerBindingTemplate, and an error, if there is any.
func (c *FakeBrokerBindingTemplates) Create(brokerBindingTemplate *experimental.BrokerBindingTemplate) (result *experimental.BrokerBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(brokerbindingtemplatesResource, brokerBindingTemplate), &experimental.BrokerBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerBindingTemplate), err
}

// Update takes the representation of a brokerBindingTemplate and updates it. Returns the server's representation of the brokerBindingTemplate, and an error, if there is any.
func (c *FakeBrokerBindingTemplates) Update(brokerBindingTemplate *experimental.BrokerBindingTemplate) (result *experimental.BrokerBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(brokerbindingtemplatesResource, brokerBindingTemplate), &experimental.BrokerBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerBindingTemplate), err
}

// Delete takes name of the brokerBindingTemplate and deletes it. Returns an error if one occurs.
func (c *FakeBrokerBindingTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(brokerbindingtemplatesResource, name), &experimental.BrokerBindingTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBrokerBindingTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(brokerbindingtemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.BrokerBindingTemplateList{})
	return err
}

// Patch applies the patch and returns the patched brokerBindingTemplate.
func (c *FakeBrokerBindingTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BrokerBindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(brokerbindingtemplatesResource, name, data, subresources...), &experimental.BrokerBindingTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerBindingTemplate), err
}
