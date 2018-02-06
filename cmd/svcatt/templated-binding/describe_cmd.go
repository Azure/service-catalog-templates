// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedbinding

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/output"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	*svcattcommand.Context
	ns       string
	name     string
	traverse bool
}

// NewDescribeCmd builds a "svcat describe templated-binding" command
func NewDescribeCmd(cxt *svcattcommand.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "templated-binding NAME",
		Aliases: []string{"templated-bindings", "templatedbindings", "templatedbinding", "tbnd"},
		Short:   "Show details of a specific templated binding",
		Example: `
  svcat describe templated-binding wordpress-mysql-binding
`,
		PreRunE: command.PreRunE(describeCmd),
		RunE:    command.RunE(describeCmd),
	}
	cmd.Flags().StringVarP(
		&describeCmd.ns,
		"namespace",
		"n",
		"",
		"The namespace in which to get the resource",
	)
	cmd.Flags().BoolVarP(
		&describeCmd.traverse,
		"traverse",
		"t",
		false,
		"Whether or not to traverse from binding -> instance -> class/plan -> broker",
	)
	return cmd
}

func (c *describeCmd) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("name is required")
	}
	c.name = args[0]

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	return nil
}

func (c *describeCmd) Run() error {
	return c.describe()
}

func (c *describeCmd) describe() error {
	tbnd, err := c.App().RetrieveTemplatedBinding(c.ns, c.name)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedBindingDetails(c.Output, tbnd)

	if c.traverse {
		tinst, class, plan, broker, err := c.App().TemplatedBindingParentHierarchy(tbnd)
		if err != nil {
			return fmt.Errorf("unable to traverse up the templated binding hierarchy (%s)", err)
		}
		svcattoutput.WriteParentTemplatedInstance(c.Output, tinst)
		output.WriteParentClass(c.Output, class)
		output.WriteParentPlan(c.Output, plan)
		output.WriteParentBroker(c.Output, broker)
	}

	return nil
}
