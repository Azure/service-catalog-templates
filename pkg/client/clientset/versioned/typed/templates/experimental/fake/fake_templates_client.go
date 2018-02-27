package fake

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/typed/templates/experimental"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeTemplatesExperimental struct {
	*testing.Fake
}

func (c *FakeTemplatesExperimental) BindingTemplates(namespace string) experimental.BindingTemplateInterface {
	return &FakeBindingTemplates{c, namespace}
}

func (c *FakeTemplatesExperimental) InstanceTemplates(namespace string) experimental.InstanceTemplateInterface {
	return &FakeInstanceTemplates{c, namespace}
}

func (c *FakeTemplatesExperimental) TemplatedBindings(namespace string) experimental.TemplatedBindingInterface {
	return &FakeTemplatedBindings{c, namespace}
}

func (c *FakeTemplatesExperimental) TemplatedInstances(namespace string) experimental.TemplatedInstanceInterface {
	return &FakeTemplatedInstances{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeTemplatesExperimental) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
