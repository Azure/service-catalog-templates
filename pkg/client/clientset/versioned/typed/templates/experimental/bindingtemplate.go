package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	scheme "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BindingTemplatesGetter has a method to return a BindingTemplateInterface.
// A group's client should implement this interface.
type BindingTemplatesGetter interface {
	BindingTemplates(namespace string) BindingTemplateInterface
}

// BindingTemplateInterface has methods to work with BindingTemplate resources.
type BindingTemplateInterface interface {
	Create(*experimental.BindingTemplate) (*experimental.BindingTemplate, error)
	Update(*experimental.BindingTemplate) (*experimental.BindingTemplate, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*experimental.BindingTemplate, error)
	List(opts v1.ListOptions) (*experimental.BindingTemplateList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BindingTemplate, err error)
	BindingTemplateExpansion
}

// bindingTemplates implements BindingTemplateInterface
type bindingTemplates struct {
	client rest.Interface
	ns     string
}

// newBindingTemplates returns a BindingTemplates
func newBindingTemplates(c *TemplatesExperimentalClient, namespace string) *bindingTemplates {
	return &bindingTemplates{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the bindingTemplate, and returns the corresponding bindingTemplate object, and an error if there is any.
func (c *bindingTemplates) Get(name string, options v1.GetOptions) (result *experimental.BindingTemplate, err error) {
	result = &experimental.BindingTemplate{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bindingtemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BindingTemplates that match those selectors.
func (c *bindingTemplates) List(opts v1.ListOptions) (result *experimental.BindingTemplateList, err error) {
	result = &experimental.BindingTemplateList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bindingtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bindingTemplates.
func (c *bindingTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bindingtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a bindingTemplate and creates it.  Returns the server's representation of the bindingTemplate, and an error, if there is any.
func (c *bindingTemplates) Create(bindingTemplate *experimental.BindingTemplate) (result *experimental.BindingTemplate, err error) {
	result = &experimental.BindingTemplate{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bindingtemplates").
		Body(bindingTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a bindingTemplate and updates it. Returns the server's representation of the bindingTemplate, and an error, if there is any.
func (c *bindingTemplates) Update(bindingTemplate *experimental.BindingTemplate) (result *experimental.BindingTemplate, err error) {
	result = &experimental.BindingTemplate{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bindingtemplates").
		Name(bindingTemplate.Name).
		Body(bindingTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the bindingTemplate and deletes it. Returns an error if one occurs.
func (c *bindingTemplates) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bindingtemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bindingTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bindingtemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched bindingTemplate.
func (c *bindingTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *experimental.BindingTemplate, err error) {
	result = &experimental.BindingTemplate{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bindingtemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
