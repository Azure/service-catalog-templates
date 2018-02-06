// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcattoutput

import (
	"io"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/output"
)

// WriteBindingTemplateList prints a list of binding templates.
func WriteBindingTemplateList(w io.Writer, bndts ...templates.BindingTemplateInterface) {
	t := output.NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Scope",
		"Service Type",
	})

	for _, bndt := range bndts {

		t.Append([]string{
			bndt.GetName(),
			getScopeText(bndt.GetScope(), bndt.GetScopeName()),
			bndt.GetServiceType(),
		})
	}

	t.Render()
}

// WriteTemplatedBindingDetails prints an binding template.
func WriteBindingTemplateDetails(w io.Writer, bndt templates.BindingTemplateInterface) {
	t := output.NewDetailsTable(w)

	t.AppendBulk([][]string{
		{"Name:", bndt.GetName()},
		{"Scope:", getScopeText(bndt.GetScope(), bndt.GetScopeName())},
		{"Service Type:", bndt.GetServiceType()},
	})

	t.Render()
}
