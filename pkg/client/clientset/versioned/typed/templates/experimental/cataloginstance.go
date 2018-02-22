package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	scheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CatalogInstancesGetter has a method to return a CatalogInstanceInterface.
// A group's client should implement this interface.
type CatalogInstancesGetter interface {
	CatalogInstances(namespace string) CatalogInstanceInterface
}

// CatalogInstanceInterface has methods to work with CatalogInstance resources.
type CatalogInstanceInterface interface {
	Create(*experimental.CatalogInstance) (*experimental.CatalogInstance, error)
	Update(*experimental.CatalogInstance) (*experimental.CatalogInstance, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.CatalogInstance, error)
	List(opts v1.ListOptions) (*experimental.CatalogInstanceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogInstance, err error)
	CatalogInstanceExpansion
}

// catalogInstances implements CatalogInstanceInterface
type catalogInstances struct {
	client rest.Interface
	ns     string
}

// newCatalogInstances returns a CatalogInstances
func newCatalogInstances(c *TemplatesExperimentalClient, namespace string) *catalogInstances {
	return &catalogInstances{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the catalogInstance, and returns the corresponding catalogInstance object, and an error if there is any.
func (c *catalogInstances) Get(name string, options v1.GetOptions) (result *experimental.CatalogInstance, err error) {
	result = &experimental.CatalogInstance{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cataloginstances").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CatalogInstances that match those selectors.
func (c *catalogInstances) List(opts v1.ListOptions) (result *experimental.CatalogInstanceList, err error) {
	result = &experimental.CatalogInstanceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cataloginstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested catalogInstances.
func (c *catalogInstances) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cataloginstances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a catalogInstance and creates it.  Returns the server's representation of the catalogInstance, and an error, if there is any.
func (c *catalogInstances) Create(catalogInstance *experimental.CatalogInstance) (result *experimental.CatalogInstance, err error) {
	result = &experimental.CatalogInstance{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cataloginstances").
		Body(catalogInstance).
		Do().
		Into(result)
	return
}

// Update takes the representation of a catalogInstance and updates it. Returns the server's representation of the catalogInstance, and an error, if there is any.
func (c *catalogInstances) Update(catalogInstance *experimental.CatalogInstance) (result *experimental.CatalogInstance, err error) {
	result = &experimental.CatalogInstance{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cataloginstances").
		Name(catalogInstance.Name).
		Body(catalogInstance).
		Do().
		Into(result)
	return
}

// Delete takes name of the catalogInstance and deletes it. Returns an error if one occurs.
func (c *catalogInstances) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cataloginstances").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *catalogInstances) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cataloginstances").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched catalogInstance.
func (c *catalogInstances) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogInstance, err error) {
	result = &experimental.CatalogInstance{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cataloginstances").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
