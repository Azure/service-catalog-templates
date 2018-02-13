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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	ClusterServiceClassExternalName string `json:"clusterServiceClassExternalName"`
	ClusterServicePlanExternalName  string `json:"clusterServicePlanExternalName"`
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

// Instance is a specification for a Instance resource
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status"`
}

// InstanceSpec is the spec for a Instance resource
type InstanceSpec struct {
	PlanSelector      metav1.LabelSelector `json:"planSelector,omitempty"`
	ClassExternalName string               `json:"classExternalName,omitempty"`
	PlanExternalName  string               `json:"planExternalName,omitempty"`
}

// InstanceStatus is the status for a Instance resource
type InstanceStatus struct {
	ResolvedClass ObjectReference `json:"resolvedClass"`
	ResolvedPlan  ObjectReference `json:"resolvedPlan"`
	// TODO: parameters
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceList is a list of Instance resources
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Instance `json:"items"`
}

// ClusterObjectReference contains enough information to let you locate the
// cluster-scoped referenced object.
type ObjectReference struct {
	// Name of the referent.
	Name string
	// Namespace of the referent. Optional for cluster scoped resources.
	//Namespace string
}