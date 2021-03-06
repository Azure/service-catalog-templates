// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

// This file was automatically generated by lister-gen

package experimental

import (
	experimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// InstanceTemplateLister helps list InstanceTemplates.
type InstanceTemplateLister interface {
	// List lists all InstanceTemplates in the indexer.
	List(selector labels.Selector) (ret []*experimental.InstanceTemplate, err error)
	// InstanceTemplates returns an object that can list and get InstanceTemplates.
	InstanceTemplates(namespace string) InstanceTemplateNamespaceLister
	InstanceTemplateListerExpansion
}

// instanceTemplateLister implements the InstanceTemplateLister interface.
type instanceTemplateLister struct {
	indexer cache.Indexer
}

// NewInstanceTemplateLister returns a new InstanceTemplateLister.
func NewInstanceTemplateLister(indexer cache.Indexer) InstanceTemplateLister {
	return &instanceTemplateLister{indexer: indexer}
}

// List lists all InstanceTemplates in the indexer.
func (s *instanceTemplateLister) List(selector labels.Selector) (ret []*experimental.InstanceTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.InstanceTemplate))
	})
	return ret, err
}

// InstanceTemplates returns an object that can list and get InstanceTemplates.
func (s *instanceTemplateLister) InstanceTemplates(namespace string) InstanceTemplateNamespaceLister {
	return instanceTemplateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// InstanceTemplateNamespaceLister helps list and get InstanceTemplates.
type InstanceTemplateNamespaceLister interface {
	// List lists all InstanceTemplates in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*experimental.InstanceTemplate, err error)
	// Get retrieves the InstanceTemplate from the indexer for a given namespace and name.
	Get(name string) (*experimental.InstanceTemplate, error)
	InstanceTemplateNamespaceListerExpansion
}

// instanceTemplateNamespaceLister implements the InstanceTemplateNamespaceLister
// interface.
type instanceTemplateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all InstanceTemplates in the indexer for a given namespace.
func (s instanceTemplateNamespaceLister) List(selector labels.Selector) (ret []*experimental.InstanceTemplate, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*experimental.InstanceTemplate))
	})
	return ret, err
}

// Get retrieves the InstanceTemplate from the indexer for a given namespace and name.
func (s instanceTemplateNamespaceLister) Get(name string) (*experimental.InstanceTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(experimental.Resource("instancetemplate"), name)
	}
	return obj.(*experimental.InstanceTemplate), nil
}
