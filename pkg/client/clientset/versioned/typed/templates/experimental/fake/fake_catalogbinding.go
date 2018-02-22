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

// FakeCatalogBindings implements CatalogBindingInterface
type FakeCatalogBindings struct {
	Fake *FakeTemplatesExperimental
	ns   string
}

var catalogbindingsResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "catalogbindings"}

var catalogbindingsKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "CatalogBinding"}

// Get takes name of the catalogBinding, and returns the corresponding catalogBinding object, and an error if there is any.
func (c *FakeCatalogBindings) Get(name string, options v1.GetOptions) (result *experimental.CatalogBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(catalogbindingsResource, c.ns, name), &experimental.CatalogBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogBinding), err
}

// List takes label and field selectors, and returns the list of CatalogBindings that match those selectors.
func (c *FakeCatalogBindings) List(opts v1.ListOptions) (result *experimental.CatalogBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(catalogbindingsResource, catalogbindingsKind, c.ns, opts), &experimental.CatalogBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.CatalogBindingList{}
	for _, item := range obj.(*experimental.CatalogBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested catalogBindings.
func (c *FakeCatalogBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(catalogbindingsResource, c.ns, opts))

}

// Create takes the representation of a catalogBinding and creates it.  Returns the server's representation of the catalogBinding, and an error, if there is any.
func (c *FakeCatalogBindings) Create(catalogBinding *experimental.CatalogBinding) (result *experimental.CatalogBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(catalogbindingsResource, c.ns, catalogBinding), &experimental.CatalogBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogBinding), err
}

// Update takes the representation of a catalogBinding and updates it. Returns the server's representation of the catalogBinding, and an error, if there is any.
func (c *FakeCatalogBindings) Update(catalogBinding *experimental.CatalogBinding) (result *experimental.CatalogBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(catalogbindingsResource, c.ns, catalogBinding), &experimental.CatalogBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogBinding), err
}

// Delete takes name of the catalogBinding and deletes it. Returns an error if one occurs.
func (c *FakeCatalogBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(catalogbindingsResource, c.ns, name), &experimental.CatalogBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCatalogBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(catalogbindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.CatalogBindingList{})
	return err
}

// Patch applies the patch and returns the patched catalogBinding.
func (c *FakeCatalogBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(catalogbindingsResource, c.ns, name, data, subresources...), &experimental.CatalogBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogBinding), err
}
