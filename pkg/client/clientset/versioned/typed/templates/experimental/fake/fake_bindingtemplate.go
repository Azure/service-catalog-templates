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

// FakeBindingTemplates implements BindingTemplateInterface
type FakeBindingTemplates struct {
	Fake *FakeTemplatesExperimental
	ns   string
}

var bindingtemplatesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "bindingtemplates"}

var bindingtemplatesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "BindingTemplate"}

// Get takes name of the bindingTemplate, and returns the corresponding bindingTemplate object, and an error if there is any.
func (c *FakeBindingTemplates) Get(name string, options v1.GetOptions) (result *experimental.BindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bindingtemplatesResource, c.ns, name), &experimental.BindingTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BindingTemplate), err
}

// List takes label and field selectors, and returns the list of BindingTemplates that match those selectors.
func (c *FakeBindingTemplates) List(opts v1.ListOptions) (result *experimental.BindingTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bindingtemplatesResource, bindingtemplatesKind, c.ns, opts), &experimental.BindingTemplateList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.BindingTemplateList{}
	for _, item := range obj.(*experimental.BindingTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested bindingTemplates.
func (c *FakeBindingTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(bindingtemplatesResource, c.ns, opts))

}

// Create takes the representation of a bindingTemplate and creates it.  Returns the server's representation of the bindingTemplate, and an error, if there is any.
func (c *FakeBindingTemplates) Create(bindingTemplate *experimental.BindingTemplate) (result *experimental.BindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bindingtemplatesResource, c.ns, bindingTemplate), &experimental.BindingTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BindingTemplate), err
}

// Update takes the representation of a bindingTemplate and updates it. Returns the server's representation of the bindingTemplate, and an error, if there is any.
func (c *FakeBindingTemplates) Update(bindingTemplate *experimental.BindingTemplate) (result *experimental.BindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bindingtemplatesResource, c.ns, bindingTemplate), &experimental.BindingTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BindingTemplate), err
}

// Delete takes name of the bindingTemplate and deletes it. Returns an error if one occurs.
func (c *FakeBindingTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bindingtemplatesResource, c.ns, name), &experimental.BindingTemplate{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBindingTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bindingtemplatesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.BindingTemplateList{})
	return err
}

// Patch applies the patch and returns the patched bindingTemplate.
func (c *FakeBindingTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BindingTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bindingtemplatesResource, c.ns, name, data, subresources...), &experimental.BindingTemplate{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BindingTemplate), err
}
