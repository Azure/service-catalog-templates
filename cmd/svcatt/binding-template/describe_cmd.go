// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package bindingtemplate

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/spf13/cobra"
)

type describeCmd struct {
	*svcattcommand.Context
	ns           string
	name         string
	traverse     bool
	brokerLevel  bool
	clusterLevel bool
	serviceType  string
}

// NewDescribeCmd builds a "svcat describe binding-template" command
func NewDescribeCmd(cxt *svcattcommand.Context) *cobra.Command {
	describeCmd := &describeCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "binding-template NAME",
		Aliases: []string{"binding-template", "bindingtemplates", "bindingtemplate", "bndt"},
		Short:   "Show details of a specific binding template",
		Example: `
  svcat describe binding-template default-mysqldb
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
	cmd.Flags().BoolVarP(&describeCmd.brokerLevel, "broker", "b", false,
		"List templates defined at the broker-level")
	cmd.Flags().BoolVarP(&describeCmd.clusterLevel, "cluster", "c", false,
		"List templates defined at the cluster-level")
	cmd.Flags().StringVarP(&describeCmd.serviceType, "type", "t", "",
		"Filter the templates by a service type")
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
	bndt, err := c.App().GetBindingTemplate(c.ns, c.clusterLevel, c.brokerLevel, c.name)
	if err != nil {
		return err
	}

	svcattoutput.WriteBindingTemplateDetails(c.Output, bndt)

	return nil
}
