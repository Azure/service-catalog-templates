package experimental

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type BindingTemplateInterface interface {
	GetParameters() *runtime.RawExtension
	GetParametersFrom() []svcat.ParametersFromSource
	GetSecretKeys() map[string]string
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

func (t *ClusterBindingTemplate) GetParameters() *runtime.RawExtension {
	return t.Spec.Parameters
}

func (t *ClusterBindingTemplate) GetParametersFrom() []svcat.ParametersFromSource {
	return t.Spec.ParametersFrom
}

func (t *ClusterBindingTemplate) GetSecretKeys() map[string]string {
	return t.Spec.SecretKeys
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
