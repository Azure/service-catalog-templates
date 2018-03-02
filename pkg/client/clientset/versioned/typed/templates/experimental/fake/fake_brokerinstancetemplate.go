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

// FakeBrokerInstanceTemplates implements BrokerInstanceTemplateInterface
type FakeBrokerInstanceTemplates struct {
	Fake *FakeTemplatesExperimental
}

var brokerinstancetemplatesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "brokerinstancetemplates"}

var brokerinstancetemplatesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "BrokerInstanceTemplate"}

// Get takes name of the brokerInstanceTemplate, and returns the corresponding brokerInstanceTemplate object, and an error if there is any.
func (c *FakeBrokerInstanceTemplates) Get(name string, options v1.GetOptions) (result *experimental.BrokerInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(brokerinstancetemplatesResource, name), &experimental.BrokerInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerInstanceTemplate), err
}

// List takes label and field selectors, and returns the list of BrokerInstanceTemplates that match those selectors.
func (c *FakeBrokerInstanceTemplates) List(opts v1.ListOptions) (result *experimental.BrokerInstanceTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(brokerinstancetemplatesResource, brokerinstancetemplatesKind, opts), &experimental.BrokerInstanceTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.BrokerInstanceTemplateList{}
	for _, item := range obj.(*experimental.BrokerInstanceTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested brokerInstanceTemplates.
func (c *FakeBrokerInstanceTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(brokerinstancetemplatesResource, opts))
}

// Create takes the representation of a brokerInstanceTemplate and creates it.  Returns the server's representation of the brokerInstanceTemplate, and an error, if there is any.
func (c *FakeBrokerInstanceTemplates) Create(brokerInstanceTemplate *experimental.BrokerInstanceTemplate) (result *experimental.BrokerInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(brokerinstancetemplatesResource, brokerInstanceTemplate), &experimental.BrokerInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerInstanceTemplate), err
}

// Update takes the representation of a brokerInstanceTemplate and updates it. Returns the server's representation of the brokerInstanceTemplate, and an error, if there is any.
func (c *FakeBrokerInstanceTemplates) Update(brokerInstanceTemplate *experimental.BrokerInstanceTemplate) (result *experimental.BrokerInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(brokerinstancetemplatesResource, brokerInstanceTemplate), &experimental.BrokerInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerInstanceTemplate), err
}

// Delete takes name of the brokerInstanceTemplate and deletes it. Returns an error if one occurs.
func (c *FakeBrokerInstanceTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(brokerinstancetemplatesResource, name), &experimental.BrokerInstanceTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBrokerInstanceTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(brokerinstancetemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.BrokerInstanceTemplateList{})
	return err
}

// Patch applies the patch and returns the patched brokerInstanceTemplate.
func (c *FakeBrokerInstanceTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BrokerInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(brokerinstancetemplatesResource, name, data, subresources...), &experimental.BrokerInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.BrokerInstanceTemplate), err
}
