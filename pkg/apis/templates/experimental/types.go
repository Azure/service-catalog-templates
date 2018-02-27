/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package experimental

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

const (
	FieldServiceTypeName = "serviceType"
)

var (
	InstanceKind = strings.Split(fmt.Sprintf("%T", TemplatedInstance{}), ".")[1]
	BindingKind  = strings.Split(fmt.Sprintf("%T", TemplatedBinding{}), ".")[1]
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceTemplate is a specification for a InstanceTemplate resource
type InstanceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceTemplateSpec   `json:"spec"`
	Status InstanceTemplateStatus `json:"status"`
}

// InstanceTemplateSpec is the spec for a InstanceTemplate resource
type InstanceTemplateSpec struct {
	ServiceType string `json:"serviceType"`

	// TODO: Should this switch to using servicecatalog's PlanReference type?
	ClassExternalName string `json:"classExternalName"`
	PlanExternalName  string `json:"planExternalName"`

	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// +optional
	ParametersFrom []svcat.ParametersFromSource `json:"parametersFrom,omitempty"`
}

// InstanceTemplateStatus is the status for a InstanceTemplate resource
type InstanceTemplateStatus struct {
	Message int32 `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceTemplateList is a list of InstanceTemplate resources
type InstanceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []InstanceTemplate `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplatedInstance is a specification for a TemplatedInstance resource
type TemplatedInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplatedInstanceSpec   `json:"spec"`
	Status TemplatedInstanceStatus `json:"status"`
}

// TemplatedInstanceSpec is the spec for a TemplatedInstance resource
type TemplatedInstanceSpec struct {
	ServiceType string `json:"serviceType"`

	// +optional
	PlanSelector *metav1.LabelSelector `json:"planSelector,omitempty"`

	// TODO: Support all the same fields as ServiceInstanceSpec: external name, external id, k8s name (uuid)

	// +optional
	ClassExternalName string `json:"classExternalName,omitempty"`

	// +optional
	PlanExternalName string `json:"planExternalName,omitempty"`

	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// +optional
	ParametersFrom []svcat.ParametersFromSource `json:"parametersFrom,omitempty"`

	// Immutable.
	// +optional
	ExternalID string `json:"externalID"`

	// +optional
	UpdateRequests int64 `json:"updateRequests"`
}

// TemplatedInstanceStatus is the status for a TemplatedInstance resource
type TemplatedInstanceStatus struct {
	ResolvedClass svcat.ObjectReference `json:"resolvedClass"`
	ResolvedPlan  svcat.ObjectReference `json:"resolvedPlan"`
	// TODO: parameters
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplatedInstanceList is a list of TemplatedInstance resources
type TemplatedInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []TemplatedInstance `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BindingTemplate is a specification for a BindingTemplate resource
type BindingTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BindingTemplateSpec   `json:"spec"`
	Status BindingTemplateStatus `json:"status"`
}

// BindingTemplateSpec is the spec for a BindingTemplate resource
type BindingTemplateSpec struct {
	ServiceType string `json:"serviceType"`

	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// +optional
	ParametersFrom []svcat.ParametersFromSource `json:"parametersFrom,omitempty"`

	// +optional
	SecretKeys map[string]string `json:"secretKeys,omitempty"`
}

// BindingTemplateStatus is the status for a BindingTemplate resource
type BindingTemplateStatus struct {
	Message int32 `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BindingTemplateList is a list of BindingTemplate resources
type BindingTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []BindingTemplate `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplatedBinding is a specification for a TemplatedBinding resource
type TemplatedBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplatedBindingSpec   `json:"spec"`
	Status TemplatedBindingStatus `json:"status"`
}

// TemplatedBindingSpec is the spec for a TemplatedBinding resource
type TemplatedBindingSpec struct {
	// Immutable.
	InstanceRef svcat.LocalObjectReference `json:"instanceRef"`

	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// +optional
	ParametersFrom []svcat.ParametersFromSource `json:"parametersFrom,omitempty"`

	// +optional
	SecretKeys map[string]string `json:"secretKeys,omitempty"`

	SecretName string `json:"secretName,omitempty"`

	// Immutable.
	// +optional
	ExternalID string `json:"externalID,omitempty"`
}

// TemplatedBindingStatus is the status for a TemplatedBinding resource
type TemplatedBindingStatus struct {
	// TODO: parameters, secretKeys
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplatedBindingList is a list of TemplatedBinding resources
type TemplatedBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []TemplatedBinding `json:"items"`
}
