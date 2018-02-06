// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedinstance

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/spf13/cobra"
)

type deprovisonCmd struct {
	*svcattcommand.Context
	ns           string
	instanceName string
}

// NewDeprovisionCmd builds a "svcat deprovision" command
func NewDeprovisionCmd(cxt *svcattcommand.Context) *cobra.Command {
	deprovisonCmd := &deprovisonCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "deprovision NAME",
		Short: "Deletes an instance of a service",
		Example: `
  svcat deprovision wordpress-mysql-instance
`,
		PreRunE: command.PreRunE(deprovisonCmd),
		RunE:    command.RunE(deprovisonCmd),
	}
	cmd.Flags().StringVarP(&deprovisonCmd.ns, "namespace", "n", "",
		"The namespace of the resource")
	return cmd
}

func (c *deprovisonCmd) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}
	c.instanceName = args[0]

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	return nil
}

func (c *deprovisonCmd) Run() error {
	return c.deprovision()
}

func (c *deprovisonCmd) deprovision() error {
	return c.App().Deprovision(c.ns, c.instanceName)
}
