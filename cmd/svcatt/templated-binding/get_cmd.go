// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedbinding

import (
	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/spf13/cobra"
)

type getCmd struct {
	*svcattcommand.Context
	ns            string
	name          string
	allNamespaces bool
}

// NewGetCmd builds a "svcat get templated-bindings" command
func NewGetCmd(cxt *svcattcommand.Context) *cobra.Command {
	getCmd := &getCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "templated-bindings [name]",
		Aliases: []string{"templated-binding", "tempaltedbindings", "templatedbinding", "tbnd"},
		Short:   "List templated bindings, optionally filtered by name",
		Example: `
  svcat get templated-bindings
  svcat get templated-bindings --all-namespaces
  svcat get templated-binding wordpress-mysql-binding
  svcat get templated-binding -n ci concourse-postgres-binding
`,
		PreRunE: command.PreRunE(getCmd),
		RunE:    command.RunE(getCmd),
	}

	cmd.Flags().StringVarP(
		&getCmd.ns,
		"namespace",
		"n",
		"",
		"The namespace from which to get the resources",
	)
	cmd.Flags().BoolVarP(
		&getCmd.allNamespaces,
		"all-namespaces",
		"",
		false,
		"List all bindings across namespaces",
	)
	return cmd
}

func (c *getCmd) Validate(args []string) error {
	if len(args) > 0 {
		c.name = args[0]
	}

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	return nil
}

func (c *getCmd) Run() error {
	if c.name == "" {
		return c.getAll()
	}

	return c.get()
}

func (c *getCmd) getAll() error {
	if c.allNamespaces {
		c.ns = ""
	}

	tbnds, err := c.App().RetrieveTemplatedBindings(c.ns)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedBindingList(c.Output, tbnds.Items...)
	return nil
}

func (c *getCmd) get() error {
	tbnd, err := c.App().RetrieveTemplatedBinding(c.ns, c.name)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedBindingList(c.Output, *tbnd)
	return nil
}
