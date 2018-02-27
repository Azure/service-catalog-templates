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

// FakeTemplatedInstances implements TemplatedInstanceInterface
type FakeTemplatedInstances struct {
	Fake *FakeTemplatesExperimental
	ns   string
}

var templatedinstancesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "templatedinstances"}

var templatedinstancesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "TemplatedInstance"}

// Get takes name of the templatedInstance, and returns the corresponding templatedInstance object, and an error if there is any.
func (c *FakeTemplatedInstances) Get(name string, options v1.GetOptions) (result *experimental.TemplatedInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(templatedinstancesResource, c.ns, name), &experimental.TemplatedInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedInstance), err
}

// List takes label and field selectors, and returns the list of TemplatedInstances that match those selectors.
func (c *FakeTemplatedInstances) List(opts v1.ListOptions) (result *experimental.TemplatedInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(templatedinstancesResource, templatedinstancesKind, c.ns, opts), &experimental.TemplatedInstanceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.TemplatedInstanceList{}
	for _, item := range obj.(*experimental.TemplatedInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested templatedInstances.
func (c *FakeTemplatedInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(templatedinstancesResource, c.ns, opts))

}

// Create takes the representation of a templatedInstance and creates it.  Returns the server's representation of the templatedInstance, and an error, if there is any.
func (c *FakeTemplatedInstances) Create(templatedInstance *experimental.TemplatedInstance) (result *experimental.TemplatedInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(templatedinstancesResource, c.ns, templatedInstance), &experimental.TemplatedInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedInstance), err
}

// Update takes the representation of a templatedInstance and updates it. Returns the server's representation of the templatedInstance, and an error, if there is any.
func (c *FakeTemplatedInstances) Update(templatedInstance *experimental.TemplatedInstance) (result *experimental.TemplatedInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(templatedinstancesResource, c.ns, templatedInstance), &experimental.TemplatedInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedInstance), err
}

// Delete takes name of the templatedInstance and deletes it. Returns an error if one occurs.
func (c *FakeTemplatedInstances) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(templatedinstancesResource, c.ns, name), &experimental.TemplatedInstance{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTemplatedInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(templatedinstancesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.TemplatedInstanceList{})
	return err
}

// Patch applies the patch and returns the patched templatedInstance.
func (c *FakeTemplatedInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(templatedinstancesResource, c.ns, name, data, subresources...), &experimental.TemplatedInstance{})

	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.TemplatedInstance), err
}
