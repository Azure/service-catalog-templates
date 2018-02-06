// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcattoutput

import (
	"fmt"
	"io"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/output"
)

// WriteTemplatedInstanceList prints a list of templated instances.
func WriteTemplatedInstanceList(w io.Writer, tinsts ...templates.TemplatedInstance) {
	t := output.NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Service Type",
		"Class",
		"Plan",
		"Status",
	})

	for _, tinst := range tinsts {
		t.Append([]string{
			tinst.Name,
			tinst.Namespace,
			tinst.Spec.ServiceType,
			tinst.Spec.ClusterServiceClassExternalName,
			tinst.Spec.ClusterServicePlanExternalName,
			"",
		})
	}

	t.Render()
}

// WriteTemplatedInstanceDetails prints a templated instance.
func WriteTemplatedInstanceDetails(w io.Writer, tinst *templates.TemplatedInstance) {
	t := output.NewDetailsTable(w)

	t.AppendBulk([][]string{
		{"Name:", tinst.Name},
		{"Namespace:", tinst.Namespace},
		{"Status:", ""},
		{"Service Type:", tinst.Spec.ServiceType},
		{"Class:", tinst.Spec.ClusterServiceClassExternalName},
		{"Plan:", tinst.Spec.ClusterServicePlanExternalName},
	})

	t.Render()
}

// WriteParentInstance prints identifying information for a parent instance.
func WriteParentTemplatedInstance(w io.Writer, tinst *templates.TemplatedInstance) {
	fmt.Fprintln(w, "\nInstance:")
	t := output.NewDetailsTable(w)
	t.AppendBulk([][]string{
		{"Name:", tinst.Name},
		{"Namespace:", tinst.Namespace},
		{"Status:", ""},
	})
	t.Render()
}
