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

// FakeClusterInstanceTemplates implements ClusterInstanceTemplateInterface
type FakeClusterInstanceTemplates struct {
	Fake *FakeTemplatesExperimental
}

var clusterinstancetemplatesResource = schema.GroupVersionResource{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Resource: "clusterinstancetemplates"}

var clusterinstancetemplatesKind = schema.GroupVersionKind{Group: "templates.servicecatalog.k8s.io", Version: "experimental", Kind: "ClusterInstanceTemplate"}

// Get takes name of the clusterInstanceTemplate, and returns the corresponding clusterInstanceTemplate object, and an error if there is any.
func (c *FakeClusterInstanceTemplates) Get(name string, options v1.GetOptions) (result *experimental.ClusterInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(clusterinstancetemplatesResource, name), &experimental.ClusterInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterInstanceTemplate), err
}

// List takes label and field selectors, and returns the list of ClusterInstanceTemplates that match those selectors.
func (c *FakeClusterInstanceTemplates) List(opts v1.ListOptions) (result *experimental.ClusterInstanceTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(clusterinstancetemplatesResource, clusterinstancetemplatesKind, opts), &experimental.ClusterInstanceTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &experimental.ClusterInstanceTemplateList{}
	for _, item := range obj.(*experimental.ClusterInstanceTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterInstanceTemplates.
func (c *FakeClusterInstanceTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(clusterinstancetemplatesResource, opts))
}

// Create takes the representation of a clusterInstanceTemplate and creates it.  Returns the server's representation of the clusterInstanceTemplate, and an error, if there is any.
func (c *FakeClusterInstanceTemplates) Create(clusterInstanceTemplate *experimental.ClusterInstanceTemplate) (result *experimental.ClusterInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(clusterinstancetemplatesResource, clusterInstanceTemplate), &experimental.ClusterInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterInstanceTemplate), err
}

// Update takes the representation of a clusterInstanceTemplate and updates it. Returns the server's representation of the clusterInstanceTemplate, and an error, if there is any.
func (c *FakeClusterInstanceTemplates) Update(clusterInstanceTemplate *experimental.ClusterInstanceTemplate) (result *experimental.ClusterInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(clusterinstancetemplatesResource, clusterInstanceTemplate), &experimental.ClusterInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterInstanceTemplate), err
}

// Delete takes name of the clusterInstanceTemplate and deletes it. Returns an error if one occurs.
func (c *FakeClusterInstanceTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(clusterinstancetemplatesResource, name), &experimental.ClusterInstanceTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterInstanceTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(clusterinstancetemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &experimental.ClusterInstanceTemplateList{})
	return err
}

// Patch applies the patch and returns the patched clusterInstanceTemplate.
func (c *FakeClusterInstanceTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.ClusterInstanceTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(clusterinstancetemplatesResource, name, data, subresources...), &experimental.ClusterInstanceTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*experimental.ClusterInstanceTemplate), err
}
