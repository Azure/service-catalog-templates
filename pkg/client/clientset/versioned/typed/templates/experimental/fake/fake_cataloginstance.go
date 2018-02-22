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

// FakeCatalogInstances implements CatalogInstanceInterface
type FakeCatalogInstances struct {
	Fake *FakeTemplatesExperimental
	ns   string
}

var cataloginstancesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "cataloginstances"}

var cataloginstancesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "CatalogInstance"}

// Get takes name of the catalogInstance, and returns the corresponding catalogInstance object, and an error if there is any.
func (c *FakeCatalogInstances) Get(name string, options v1.GetOptions) (result *experimental.CatalogInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cataloginstancesResource, c.ns, name), &experimental.CatalogInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogInstance), err
}

// List takes label and field selectors, and returns the list of CatalogInstances that match those selectors.
func (c *FakeCatalogInstances) List(opts v1.ListOptions) (result *experimental.CatalogInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cataloginstancesResource, cataloginstancesKind, c.ns, opts), &experimental.CatalogInstanceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.CatalogInstanceList{}
	for _, item := range obj.(*experimental.CatalogInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested catalogInstances.
func (c *FakeCatalogInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cataloginstancesResource, c.ns, opts))

}

// Create takes the representation of a catalogInstance and creates it.  Returns the server's representation of the catalogInstance, and an error, if there is any.
func (c *FakeCatalogInstances) Create(catalogInstance *experimental.CatalogInstance) (result *experimental.CatalogInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cataloginstancesResource, c.ns, catalogInstance), &experimental.CatalogInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogInstance), err
}

// Update takes the representation of a catalogInstance and updates it. Returns the server's representation of the catalogInstance, and an error, if there is any.
func (c *FakeCatalogInstances) Update(catalogInstance *experimental.CatalogInstance) (result *experimental.CatalogInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cataloginstancesResource, c.ns, catalogInstance), &experimental.CatalogInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogInstance), err
}

// Delete takes name of the catalogInstance and deletes it. Returns an error if one occurs.
func (c *FakeCatalogInstances) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cataloginstancesResource, c.ns, name), &experimental.CatalogInstance{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCatalogInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cataloginstancesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.CatalogInstanceList{})
	return err
}

// Patch applies the patch and returns the patched catalogInstance.
func (c *FakeCatalogInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cataloginstancesResource, c.ns, name, data, subresources...), &experimental.CatalogInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.CatalogInstance), err
}
