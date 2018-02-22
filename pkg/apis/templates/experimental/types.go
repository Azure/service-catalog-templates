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
	InstanceKind = strings.Split(fmt.Sprintf("%T", CatalogInstance{}), ".")[1]
	BindingKind  = strings.Split(fmt.Sprintf("%T", CatalogBinding{}), ".")[1]
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

// CatalogInstance is a specification for a CatalogInstance resource
type CatalogInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogInstanceSpec   `json:"spec"`
	Status CatalogInstanceStatus `json:"status"`
}

// CatalogInstanceSpec is the spec for a CatalogInstance resource
type CatalogInstanceSpec struct {
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

// CatalogInstanceStatus is the status for a CatalogInstance resource
type CatalogInstanceStatus struct {
	ResolvedClass svcat.ObjectReference `json:"resolvedClass"`
	ResolvedPlan  svcat.ObjectReference `json:"resolvedPlan"`
	// TODO: parameters
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CatalogInstanceList is a list of CatalogInstance resources
type CatalogInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CatalogInstance `json:"items"`
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
	SecretKeys map[string]string `json:"secret-keys,omitempty"`
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

// CatalogBinding is a specification for a CatalogBinding resource
type CatalogBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CatalogBindingSpec   `json:"spec"`
	Status CatalogBindingStatus `json:"status"`
}

// CatalogBindingSpec is the spec for a CatalogBinding resource
type CatalogBindingSpec struct {
	// Immutable.
	InstanceRef svcat.LocalObjectReference `json:"instanceRef"`

	// +optional
	Parameters *runtime.RawExtension `json:"parameters,omitempty"`

	// +optional
	ParametersFrom []svcat.ParametersFromSource `json:"parametersFrom,omitempty"`

	// +optional
	SecretKeys map[string]string `json:"secret-keys,omitempty"`
}

// CatalogBindingStatus is the status for a CatalogBinding resource
type CatalogBindingStatus struct {
	// TODO: parameters, secret-keys
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CatalogBindingList is a list of CatalogBinding resources
type CatalogBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CatalogBinding `json:"items"`
}
