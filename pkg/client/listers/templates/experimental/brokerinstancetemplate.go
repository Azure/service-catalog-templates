// This file was automatically generated by lister-gen

package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// BrokerInstanceTemplateLister helps list BrokerInstanceTemplates.
type BrokerInstanceTemplateLister interface {
	// List lists all BrokerInstanceTemplates in the indexer.
	List(selector labels.Selector) (ret []*experimental.BrokerInstanceTemplate, err error)
	// Get retrieves the BrokerInstanceTemplate from the index for a given name.
	Get(name string) (*experimental.BrokerInstanceTemplate, error)
	BrokerInstanceTemplateListerExpansion
}

// brokerInstanceTemplateLister implements the BrokerInstanceTemplateLister interface.
type brokerInstanceTemplateLister struct {
	indexer cache.Indexer
}

// NewBrokerInstanceTemplateLister returns a new BrokerInstanceTemplateLister.
func NewBrokerInstanceTemplateLister(indexer cache.Indexer) BrokerInstanceTemplateLister {
	return &brokerInstanceTemplateLister{indexer: indexer}
}

// List lists all BrokerInstanceTemplates in the indexer.
func (s *brokerInstanceTemplateLister) List(selector labels.Selector) (ret []*experimental.BrokerInstanceTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.BrokerInstanceTemplate))
	})
	return ret, err
}

// Get retrieves the BrokerInstanceTemplate from the index for a given name.
func (s *brokerInstanceTemplateLister) Get(name string) (*experimental.BrokerInstanceTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(experimental.Resource("brokerinstancetemplate"), name)
	}
	return obj.(*experimental.BrokerInstanceTemplate), nil
}
