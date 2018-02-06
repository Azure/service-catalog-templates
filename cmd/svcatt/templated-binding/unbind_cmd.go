// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedbinding

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/spf13/cobra"
)

type unbindCmd struct {
	*svcattcommand.Context
	ns           string
	instanceName string
	bindingName  string
}

// NewUnbindCmd builds a "svcat unbind" command
func NewUnbindCmd(cxt *svcattcommand.Context) *cobra.Command {
	unbindCmd := &unbindCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "unbind INSTANCE_NAME",
		Short: "Unbinds an instance. When an instance name is specified, all of its bindings are removed, otherwise use --name to remove a specific binding",
		Example: `
  svcat unbind wordpress-mysql-instance
  svcat unbind --name wordpress-mysql-binding
`,
		PreRunE: command.PreRunE(unbindCmd),
		RunE:    command.RunE(unbindCmd),
	}

	cmd.Flags().StringVarP(
		&unbindCmd.ns,
		"namespace",
		"n",
		"",
		"The namespace of the resource",
	)
	cmd.Flags().StringVar(
		&unbindCmd.bindingName,
		"name",
		"",
		"The name of the binding to remove",
	)
	return cmd
}

func (c *unbindCmd) Validate(args []string) error {
	if len(args) == 0 {
		if c.bindingName == "" {
			return fmt.Errorf("an instance or binding name is required")
		}
	} else {
		c.instanceName = args[0]
	}

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	return nil
}

func (c *unbindCmd) Run() error {
	if c.instanceName != "" {
		return c.unbindTemplatedInstance()
	}
	return c.deleteTemplatedBinding()
}

func (c *unbindCmd) deleteTemplatedBinding() error {
	err := c.App().DeleteTemplatedBinding(c.ns, c.bindingName)
	if err == nil {
		svcattoutput.WriteDeletedTemplatedBindingName(c.Output, c.bindingName)
	}
	return err
}

func (c *unbindCmd) unbindTemplatedInstance() error {
	bindings, err := c.App().Unbind(c.ns, c.instanceName)
	svcattoutput.WriteDeletedTemplatedBindingNames(c.Output, bindings)
	return err
}
