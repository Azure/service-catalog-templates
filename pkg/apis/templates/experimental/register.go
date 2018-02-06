// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package experimental

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: templates.GroupName, Version: "experimental"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&TemplatedBinding{},
		&TemplatedBindingList{},
		&TemplatedInstance{},
		&TemplatedInstanceList{},
		&BindingTemplate{},
		&BindingTemplateList{},
		&InstanceTemplate{},
		&InstanceTemplateList{},
		&ClusterInstanceTemplate{},
		&ClusterInstanceTemplateList{},
		&BrokerInstanceTemplate{},
		&BrokerInstanceTemplateList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
