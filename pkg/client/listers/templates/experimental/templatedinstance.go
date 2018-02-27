// This file was automatically generated by lister-gen

package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TemplatedInstanceLister helps list TemplatedInstances.
type TemplatedInstanceLister interface {
	// List lists all TemplatedInstances in the indexer.
	List(selector labels.Selector) (ret []*experimental.TemplatedInstance, err error)
	// TemplatedInstances returns an object that can list and get TemplatedInstances.
	TemplatedInstances(namespace string) TemplatedInstanceNamespaceLister
	TemplatedInstanceListerExpansion
}

// templatedInstanceLister implements the TemplatedInstanceLister interface.
type templatedInstanceLister struct {
	indexer cache.Indexer
}

// NewTemplatedInstanceLister returns a new TemplatedInstanceLister.
func NewTemplatedInstanceLister(indexer cache.Indexer) TemplatedInstanceLister {
	return &templatedInstanceLister{indexer: indexer}
}

// List lists all TemplatedInstances in the indexer.
func (s *templatedInstanceLister) List(selector labels.Selector) (ret []*experimental.TemplatedInstance, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.TemplatedInstance))
	})
	return ret, err
}

// TemplatedInstances returns an object that can list and get TemplatedInstances.
func (s *templatedInstanceLister) TemplatedInstances(namespace string) TemplatedInstanceNamespaceLister {
	return templatedInstanceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TemplatedInstanceNamespaceLister helps list and get TemplatedInstances.
type TemplatedInstanceNamespaceLister interface {
	// List lists all TemplatedInstances in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*experimental.TemplatedInstance, err error)
	// Get retrieves the TemplatedInstance from the indexer for a given namespace and name.
	Get(name string) (*experimental.TemplatedInstance, error)
	TemplatedInstanceNamespaceListerExpansion
}

// templatedInstanceNamespaceLister implements the TemplatedInstanceNamespaceLister
// interface.
type templatedInstanceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TemplatedInstances in the indexer for a given namespace.
func (s templatedInstanceNamespaceLister) List(selector labels.Selector) (ret []*experimental.TemplatedInstance, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.TemplatedInstance))
	})
	return ret, err
}

// Get retrieves the TemplatedInstance from the indexer for a given namespace and name.
func (s templatedInstanceNamespaceLister) Get(name string) (*experimental.TemplatedInstance, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(experimental.Resource("templatedinstance"), name)
	}
	return obj.(*experimental.TemplatedInstance), nil
}