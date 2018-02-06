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

// FakeTemplatedBindings implements TemplatedBindingInterface
type FakeTemplatedBindings struct {
	Fake *FakeTemplatesExperimental
	ns   string
}

var templatedbindingsResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "templatedbindings"}

var templatedbindingsKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "TemplatedBinding"}

// Get takes name of the templatedBinding, and returns the corresponding templatedBinding object, and an error if there is any.
func (c *FakeTemplatedBindings) Get(name string, options v1.GetOptions) (result *experimental.TemplatedBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(templatedbindingsResource, c.ns, name), &experimental.TemplatedBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedBinding), err
}

// List takes label and field selectors, and returns the list of TemplatedBindings that match those selectors.
func (c *FakeTemplatedBindings) List(opts v1.ListOptions) (result *experimental.TemplatedBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(templatedbindingsResource, templatedbindingsKind, c.ns, opts), &experimental.TemplatedBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.TemplatedBindingList{}
	for _, item := range obj.(*experimental.TemplatedBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested templatedBindings.
func (c *FakeTemplatedBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(templatedbindingsResource, c.ns, opts))

}

// Create takes the representation of a templatedBinding and creates it.  Returns the server's representation of the templatedBinding, and an error, if there is any.
func (c *FakeTemplatedBindings) Create(templatedBinding *experimental.TemplatedBinding) (result *experimental.TemplatedBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(templatedbindingsResource, c.ns, templatedBinding), &experimental.TemplatedBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedBinding), err
}

// Update takes the representation of a templatedBinding and updates it. Returns the server's representation of the templatedBinding, and an error, if there is any.
func (c *FakeTemplatedBindings) Update(templatedBinding *experimental.TemplatedBinding) (result *experimental.TemplatedBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(templatedbindingsResource, c.ns, templatedBinding), &experimental.TemplatedBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedBinding), err
}

// Delete takes name of the templatedBinding and deletes it. Returns an error if one occurs.
func (c *FakeTemplatedBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(templatedbindingsResource, c.ns, name), &experimental.TemplatedBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTemplatedBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(templatedbindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.TemplatedBindingList{})
	return err
}

// Patch applies the patch and returns the patched templatedBinding.
func (c *FakeTemplatedBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(templatedbindingsResource, c.ns, name, data, subresources...), &experimental.TemplatedBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedBinding), err
}
