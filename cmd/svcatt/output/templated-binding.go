// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcattoutput

import (
	"fmt"
	"io"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/output"
)

// WriteTemplatedBindingList prints a list of bindings.
func WriteTemplatedBindingList(w io.Writer, bindings ...templates.TemplatedBinding) {
	t := output.NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Namespace",
		"Instance",
		"Status",
	})

	for _, binding := range bindings {
		t.Append([]string{
			binding.Name,
			binding.Namespace,
			binding.Spec.TemplatedInstanceRef.Name,
			"",
		})
	}

	t.Render()
}

// WriteTemplatedBindingDetails prints details for a single binding.
func WriteTemplatedBindingDetails(w io.Writer, binding *templates.TemplatedBinding) {
	t := output.NewDetailsTable(w)

	t.AppendBulk([][]string{
		{"Name:", binding.Name},
		{"Namespace:", binding.Namespace},
		{"Status:", ""},
		{"Instance:", binding.Spec.TemplatedInstanceRef.Name},
	})

	t.Render()
}

// WriteAssociatedTemplatedBindings prints a list of bindings associated with an instance.
func WriteAssociatedTemplatedBindings(w io.Writer, bindings []templates.TemplatedBinding) {
	fmt.Fprintln(w, "\nBindings:")
	if len(bindings) == 0 {
		fmt.Fprintln(w, "No bindings defined")
		return
	}

	t := output.NewListTable(w)
	t.SetHeader([]string{
		"Name",
		"Status",
	})
	for _, binding := range bindings {
		t.Append([]string{
			binding.Name,
			"",
		})
	}
	t.Render()
}

// WriteDeletedTemplatedBindingNames prints the names of a list of bindings
func WriteDeletedTemplatedBindingNames(w io.Writer, bindings []templates.TemplatedBinding) {
	for _, binding := range bindings {
		WriteDeletedTemplatedBindingName(w, binding.Name)
	}
}

// WriteDeletedTemplatedBindingName prints the name of a binding
func WriteDeletedTemplatedBindingName(w io.Writer, bindingName string) {
	fmt.Fprintf(w, "deleted %s\n", bindingName)
}
