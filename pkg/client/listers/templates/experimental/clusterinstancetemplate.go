// This file was automatically generated by lister-gen

package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ClusterInstanceTemplateLister helps list ClusterInstanceTemplates.
type ClusterInstanceTemplateLister interface {
	// List lists all ClusterInstanceTemplates in the indexer.
	List(selector labels.Selector) (ret []*experimental.ClusterInstanceTemplate, err error)
	// Get retrieves the ClusterInstanceTemplate from the index for a given name.
	Get(name string) (*experimental.ClusterInstanceTemplate, error)
	ClusterInstanceTemplateListerExpansion
}

// clusterInstanceTemplateLister implements the ClusterInstanceTemplateLister interface.
type clusterInstanceTemplateLister struct {
	indexer cache.Indexer
}

// NewClusterInstanceTemplateLister returns a new ClusterInstanceTemplateLister.
func NewClusterInstanceTemplateLister(indexer cache.Indexer) ClusterInstanceTemplateLister {
	return &clusterInstanceTemplateLister{indexer: indexer}
}

// List lists all ClusterInstanceTemplates in the indexer.
func (s *clusterInstanceTemplateLister) List(selector labels.Selector) (ret []*experimental.ClusterInstanceTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.ClusterInstanceTemplate))
	})
	return ret, err
}

// Get retrieves the ClusterInstanceTemplate from the index for a given name.
func (s *clusterInstanceTemplateLister) Get(name string) (*experimental.ClusterInstanceTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(experimental.Resource("clusterinstancetemplate"), name)
	}
	return obj.(*experimental.ClusterInstanceTemplate), nil
}