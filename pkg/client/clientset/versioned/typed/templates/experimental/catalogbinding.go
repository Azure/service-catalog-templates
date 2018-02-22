package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	scheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CatalogBindingsGetter has a method to return a CatalogBindingInterface.
// A group's client should implement this interface.
type CatalogBindingsGetter interface {
	CatalogBindings(namespace string) CatalogBindingInterface
}

// CatalogBindingInterface has methods to work with CatalogBinding resources.
type CatalogBindingInterface interface {
	Create(*experimental.CatalogBinding) (*experimental.CatalogBinding, error)
	Update(*experimental.CatalogBinding) (*experimental.CatalogBinding, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.CatalogBinding, error)
	List(opts v1.ListOptions) (*experimental.CatalogBindingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogBinding, err error)
	CatalogBindingExpansion
}

// catalogBindings implements CatalogBindingInterface
type catalogBindings struct {
	client rest.Interface
	ns     string
}

// newCatalogBindings returns a CatalogBindings
func newCatalogBindings(c *TemplatesExperimentalClient, namespace string) *catalogBindings {
	return &catalogBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the catalogBinding, and returns the corresponding catalogBinding object, and an error if there is any.
func (c *catalogBindings) Get(name string, options v1.GetOptions) (result *experimental.CatalogBinding, err error) {
	result = &experimental.CatalogBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("catalogbindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CatalogBindings that match those selectors.
func (c *catalogBindings) List(opts v1.ListOptions) (result *experimental.CatalogBindingList, err error) {
	result = &experimental.CatalogBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("catalogbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested catalogBindings.
func (c *catalogBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("catalogbindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a catalogBinding and creates it.  Returns the server's representation of the catalogBinding, and an error, if there is any.
func (c *catalogBindings) Create(catalogBinding *experimental.CatalogBinding) (result *experimental.CatalogBinding, err error) {
	result = &experimental.CatalogBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("catalogbindings").
		Body(catalogBinding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a catalogBinding and updates it. Returns the server's representation of the catalogBinding, and an error, if there is any.
func (c *catalogBindings) Update(catalogBinding *experimental.CatalogBinding) (result *experimental.CatalogBinding, err error) {
	result = &experimental.CatalogBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("catalogbindings").
		Name(catalogBinding.Name).
		Body(catalogBinding).
		Do().
		Into(result)
	return
}

// Delete takes name of the catalogBinding and deletes it. Returns an error if one occurs.
func (c *catalogBindings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("catalogbindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *catalogBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("catalogbindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched catalogBinding.
func (c *catalogBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.CatalogBinding, err error) {
	result = &experimental.CatalogBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("catalogbindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
