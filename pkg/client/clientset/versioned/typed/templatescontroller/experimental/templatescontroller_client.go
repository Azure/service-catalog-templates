package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templatescontroller/experimental"
	"github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type TemplatesExperimentalInterface interface {
	RESTClient() rest.Interface
	InstanceTemplatesGetter
}

// TemplatesExperimentalClient is used to interact with features provided by the templates.servicecatalog.k8s.io group.
type TemplatesExperimentalClient struct {
	restClient rest.Interface
}

func (c *TemplatesExperimentalClient) InstanceTemplates(namespace string) InstanceTemplateInterface {
	return newInstanceTemplates(c, namespace)
}

// NewForConfig creates a new TemplatesExperimentalClient for the given config.
func NewForConfig(c *rest.Config) (*TemplatesExperimentalClient, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &TemplatesExperimentalClient{client}, nil
}

// NewForConfigOrDie creates a new TemplatesExperimentalClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *TemplatesExperimentalClient {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new TemplatesExperimentalClient for the given RESTClient.
func New(c rest.Interface) *TemplatesExperimentalClient {
	return &TemplatesExperimentalClient{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := experimental.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *TemplatesExperimentalClient) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
