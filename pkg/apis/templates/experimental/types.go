// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

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
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerInstanceTemplate is a specification for a BrokerInstanceTemplate resource
type BrokerInstanceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerInstanceTemplateSpec   `json:"spec"`
	Status BrokerInstanceTemplateStatus `json:"status"`
}

// BrokerInstanceTemplateSpec is the spec for a BrokerInstanceTemplate resource
type BrokerInstanceTemplateSpec struct {
	InstanceTemplateSpec `json:",inline"`

	BrokerName string `json:"brokerName"`
}

// BrokerInstanceTemplateStatus is the status for a BrokerInstanceTemplate resource
type BrokerInstanceTemplateStatus InstanceTemplateStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerInstanceTemplateList is a list of BrokerInstanceTemplate resources
type BrokerInstanceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []BrokerInstanceTemplate `json:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterInstanceTemplate is a specification for a ClusterInstanceTemplate resource
type ClusterInstanceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterInstanceTemplateSpec   `json:"spec"`
	Status ClusterInstanceTemplateStatus `json:"status"`
}

// ClusterInstanceTemplateSpec is the spec for a ClusterInstanceTemplate resource
type ClusterInstanceTemplateSpec InstanceTemplateSpec

// ClusterInstanceTemplateStatus is the status for a ClusterInstanceTemplate resource
type ClusterInstanceTemplateStatus InstanceTemplateStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterInstanceTemplateList is a list of ClusterInstanceTemplate resources
type ClusterInstanceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterInstanceTemplate `json:"items"`
}

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
	// Specification of what ServiceClass/ServicePlan is being provisioned.
	svcat.PlanReference `json:",inline"`

	ServiceType string `json:"serviceType"`

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
	// +optional
	ServiceType string `json:"serviceType"`

	// +optional
	PlanSelector *metav1.LabelSelector `json:"planSelector,omitempty"`

	// Specification of what ServiceClass/ServicePlan is being provisioned.
	svcat.PlanReference `json:",inline"`

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
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerBindingTemplate is a specification for a BrokerBindingTemplate resource
type BrokerBindingTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BrokerBindingTemplateSpec   `json:"spec"`
	Status BrokerBindingTemplateStatus `json:"status"`
}

// BrokerBindingTemplateSpec is the spec for a BrokerBindingTemplate resource
type BrokerBindingTemplateSpec struct {
	BindingTemplateSpec `json:",inline"`

	BrokerName string `json:"brokerName"`
}

// BrokerBindingTemplateStatus is the status for a BrokerBindingTemplate resource
type BrokerBindingTemplateStatus BindingTemplateStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BrokerBindingTemplateList is a list of BrokerBindingTemplate resources
type BrokerBindingTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []BrokerBindingTemplate `json:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterBindingTemplate is a specification for a ClusterBindingTemplate resource
type ClusterBindingTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterBindingTemplateSpec   `json:"spec"`
	Status ClusterBindingTemplateStatus `json:"status"`
}

// ClusterBindingTemplateSpec is the spec for a ClusterBindingTemplate resource
type ClusterBindingTemplateSpec BindingTemplateSpec

// ClusterBindingTemplateStatus is the status for a ClusterBindingTemplate resource
type ClusterBindingTemplateStatus BindingTemplateStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterBindingTemplateList is a list of ClusterBindingTemplate resources
type ClusterBindingTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterBindingTemplate `json:"items"`
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
	TemplatedInstanceRef svcat.LocalObjectReference `json:"instanceRef"`

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
