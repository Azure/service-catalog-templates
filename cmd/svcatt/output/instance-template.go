// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcattoutput

import (
	"fmt"
	"io"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/output"
)

func getScopeText(scope templates.TemplateScope, name string) string {
	if name == "" {
		return string(scope)
	}
	return fmt.Sprintf("%s (%s)", scope, name)
}

// WriteInstanceTemplateList prints a list of instance templates.
func WriteInstanceTemplateList(w io.Writer, instts ...templates.InstanceTemplateInterface) {
	t := output.NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Scope",
		"Service Type",
		"Class",
		"Plan",
	})

	for _, instt := range instts {

		plan := instt.GetPlanReference()
		t.Append([]string{
			instt.GetName(),
			getScopeText(instt.GetScope(), instt.GetScopeName()),
			instt.GetServiceType(),
			plan.ClusterServiceClassExternalName,
			plan.ClusterServicePlanExternalName,
			"",
		})
	}

	t.Render()
}

// WriteTemplatedInstanceDetails prints an instance template.
func WriteInstanceTemplateDetails(w io.Writer, instt templates.InstanceTemplateInterface) {
	t := output.NewDetailsTable(w)

	plan := instt.GetPlanReference()

	t.AppendBulk([][]string{
		{"Name:", instt.GetName()},
		{"Scope:", getScopeText(instt.GetScope(), instt.GetScopeName())},
		{"Service Type:", instt.GetServiceType()},
		{"Class:", plan.ClusterServiceClassExternalName},
		{"Plan:", plan.ClusterServicePlanExternalName},
	})

	t.Render()
}
