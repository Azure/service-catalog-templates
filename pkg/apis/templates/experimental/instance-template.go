// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package experimental

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type TemplateScope string

const (
	ScopeNamespace = "namespace"
	ScopeCluster   = "cluster"
	ScopeBroker    = "broker"
)

type InstanceTemplateInterface interface {
	GetName() string
	GetScope() TemplateScope
	GetScopeName() string
	GetServiceType() string
	GetPlanReference() svcat.PlanReference
	SetPlanReference(reference svcat.PlanReference)
	GetParameters() *runtime.RawExtension
	GetParametersFrom() []svcat.ParametersFromSource
}

func (t *InstanceTemplate) GetName() string {
	return t.Name
}

func (t *InstanceTemplate) GetScope() TemplateScope {
	return ScopeNamespace
}

func (t *InstanceTemplate) GetScopeName() string {
	return t.Namespace
}

func (t *InstanceTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *InstanceTemplate) GetPlanReference() svcat.PlanReference {
	return t.Spec.PlanReference
}

func (t *InstanceTemplate) SetPlanReference(pr svcat.PlanReference) {
	t.Spec.PlanReference = pr
}

func (t *InstanceTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *InstanceTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *ClusterInstanceTemplate) GetName() string {
	return t.Name
}

func (t *ClusterInstanceTemplate) GetScope() TemplateScope {
	return ScopeCluster
}

func (t *ClusterInstanceTemplate) GetScopeName() string {
	return ""
}

func (t *ClusterInstanceTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *ClusterInstanceTemplate) GetPlanReference() svcat.PlanReference {
	return t.Spec.PlanReference
}

func (t *ClusterInstanceTemplate) SetPlanReference(pr svcat.PlanReference) {
	t.Spec.PlanReference = pr
}

func (t *ClusterInstanceTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *ClusterInstanceTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *BrokerInstanceTemplate) GetName() string {
	return t.Name
}

func (t *BrokerInstanceTemplate) GetScope() TemplateScope {
	return ScopeBroker
}

func (t *BrokerInstanceTemplate) GetScopeName() string {
	return t.Spec.BrokerName
}

func (t *BrokerInstanceTemplate) GetServiceType() string {
	return t.Spec.ServiceType
}

func (t *BrokerInstanceTemplate) GetPlanReference() svcat.PlanReference {
	return t.Spec.PlanReference
}

func (t *BrokerInstanceTemplate) SetPlanReference(pr svcat.PlanReference) {
	t.Spec.PlanReference = pr
}

func (t *BrokerInstanceTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *BrokerInstanceTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}
