package versioned

import (
	templatesexperimental "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned/typed/templates/experimental"
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	TemplatesExperimental() templatesexperimental.TemplatesExperimentalInterface
	// Deprecated: please explicitly pick a version if possible.
	Templates() templatesexperimental.TemplatesExperimentalInterface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	templatesExperimental *templatesexperimental.TemplatesExperimentalClient
}

// TemplatesExperimental retrieves the TemplatesExperimentalClient
func (c *Clientset) TemplatesExperimental() templatesexperimental.TemplatesExperimentalInterface {
	return c.templatesExperimental
}

// Deprecated: Templates retrieves the default version of TemplatesClient.
// Please explicitly pick a version.
func (c *Clientset) Templates() templatesexperimental.TemplatesExperimentalInterface {
	return c.templatesExperimental
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.templatesExperimental, err = templatesexperimental.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.templatesExperimental = templatesexperimental.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.templatesExperimental = templatesexperimental.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
