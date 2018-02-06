// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package experimental

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BindingTemplateInterface interface {
	GetName() string
	GetScope() TemplateScope
	GetScopeName() string
	GetServiceType() string
	GetParameters() *runtime.RawExtension
	GetParametersFrom() []svcat.ParametersFromSource
	GetSecretKeys() map[string]string
}

func (t *BindingTemplate) GetName() string {
	return t.Name
}

func (t *BindingTemplate) GetScope() TemplateScope {
	return ScopeNamespace
}

func (t *BindingTemplate) GetScopeName() string {
	return t.Namespace
}

func (t *BindingTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *BindingTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *BindingTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *BindingTemplate) GetSecretKeys() map[string]string {
	return t.Spec.SecretKeys
}

func (t *ClusterBindingTemplate) GetName() string {
	return t.Name
}

func (t *ClusterBindingTemplate) GetScope() TemplateScope {
	return ScopeCluster
}

func (t *ClusterBindingTemplate) GetScopeName() string {
	return ""
}

func (t *ClusterBindingTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *ClusterBindingTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *ClusterBindingTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *ClusterBindingTemplate) GetSecretKeys() map[string]string {
	return t.Spec.SecretKeys
}

func (t *BrokerBindingTemplate) GetName() string {
	return t.Name
}

func (t *BrokerBindingTemplate) GetScope() TemplateScope {
	return ScopeBroker
}

func (t *BrokerBindingTemplate) GetScopeName() string {
	return t.Spec.BrokerName
}

func (t *BrokerBindingTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *BrokerBindingTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *BrokerBindingTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *BrokerBindingTemplate) GetSecretKeys() map[string]string {
	return t.Spec.SecretKeys
}
