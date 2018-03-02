package experimental

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type InstanceTemplateInterface interface {
	GetPlanReference() svcat.PlanReference
	SetPlanReference(reference svcat.PlanReference)
	GetParameters() *runtime.RawExtension
	GetParametersFrom() []svcat.ParametersFromSource
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
