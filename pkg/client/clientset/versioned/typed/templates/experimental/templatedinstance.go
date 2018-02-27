package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	scheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TemplatedInstancesGetter has a method to return a TemplatedInstanceInterface.
// A group's client should implement this interface.
type TemplatedInstancesGetter interface {
	TemplatedInstances(namespace string) TemplatedInstanceInterface
}

// TemplatedInstanceInterface has methods to work with TemplatedInstance resources.
type TemplatedInstanceInterface interface {
	Create(*experimental.TemplatedInstance) (*experimental.TemplatedInstance, error)
	Update(*experimental.TemplatedInstance) (*experimental.TemplatedInstance, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.TemplatedInstance, error)
	List(opts v1.ListOptions) (*experimental.TemplatedInstanceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedInstance, err error)
	TemplatedInstanceExpansion
}

// templatedInstances implements TemplatedInstanceInterface
type templatedInstances struct {
	client rest.Interface
	ns     string
}

// newTemplatedInstances returns a TemplatedInstances
func newTemplatedInstances(c *TemplatesExperimentalClient, namespace string) *templatedInstances {
	return &templatedInstances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the templatedInstance, and returns the corresponding templatedInstance object, and an error if there is any.
func (c *templatedInstances) Get(name string, options v1.GetOptions) (result *experimental.TemplatedInstance, err error) {
	result = &experimental.TemplatedInstance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("templatedinstances").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of TemplatedInstances that match those selectors.
func (c *templatedInstances) List(opts v1.ListOptions) (result *experimental.TemplatedInstanceList, err error) {
	result = &experimental.TemplatedInstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("templatedinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested templatedInstances.
func (c *templatedInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("templatedinstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a templatedInstance and creates it.  Returns the server's representation of the templatedInstance, and an error, if there is any.
func (c *templatedInstances) Create(templatedInstance *experimental.TemplatedInstance) (result *experimental.TemplatedInstance, err error) {
	result = &experimental.TemplatedInstance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("templatedinstances").
		Body(templatedInstance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a templatedInstance and updates it. Returns the server's representation of the templatedInstance, and an error, if there is any.
func (c *templatedInstances) Update(templatedInstance *experimental.TemplatedInstance) (result *experimental.TemplatedInstance, err error) {
	result = &experimental.TemplatedInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("templatedinstances").
		Name(templatedInstance.Name).
		Body(templatedInstance).
		Do().
		Into(result)
	return
}

// Delete takes name of the templatedInstance and deletes it. Returns an error if one occurs.
func (c *templatedInstances) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("templatedinstances").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *templatedInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("templatedinstances").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched templatedInstance.
func (c *templatedInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.TemplatedInstance, err error) {
	result = &experimental.TemplatedInstance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("templatedinstances").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
