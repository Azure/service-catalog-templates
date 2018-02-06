// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedinstance

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

// NewGetCmd builds a "svcat get templated-instances" command
func NewGetCmd(cxt *svcattcommand.Context) *cobra.Command {
	getCmd := &getCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "templated-instances [name]",
		Aliases: []string{"templatedinstances", "templatedinstance", "tinst"},
		Short:   "List templated instances, optionally filtered by name",
		Example: `
  svcat get templated-instances
  svcat get templated-instances --all-namespaces
  svcat get templated-instances wordpress-mysql-instance
  svcat get templated-instances -n ci concourse-postgres-instance
`,
		PreRunE: command.PreRunE(getCmd),
		RunE:    command.RunE(getCmd),
	}
	cmd.Flags().StringVarP(
		&getCmd.ns,
		"namespace",
		"n",
		"",
		"The namespace in which to get the resources",
	)
	cmd.Flags().BoolVarP(
		&getCmd.allNamespaces,
		"all-namespaces",
		"",
		false,
		"List all resources across namespaces",
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

	tinsts, err := c.App().RetrieveTemplatedInstances(c.ns)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedInstanceList(c.Output, tinsts.Items...)
	return nil
}

func (c *getCmd) get() error {
	tinst, err := c.App().RetrieveTemplatedInstance(c.ns, c.name)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedInstanceList(c.Output, *tinst)
	return nil
}
