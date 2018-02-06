// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package instancetemplate

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
	brokerLevel   bool
	clusterLevel  bool
	serviceType   string
}

// NewGetCmd builds a "svcat get instance-templates" command
func NewGetCmd(cxt *svcattcommand.Context) *cobra.Command {
	getCmd := &getCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:     "instance-templates [name]",
		Aliases: []string{"instance-template", "instancetemplates", "instancetemplate", "instt"},
		Short:   "List instance templates, optionally filtered by name and scope",
		Example: `
  svcat get instance-templates --namespace teamA
  svcat get instance-templates --cluster
  svcat get instance-templates --broker
  svcat get instance-templates --type mysqldb
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
	cmd.Flags().BoolVarP(&getCmd.brokerLevel, "broker", "b", false, "List templates defined at the broker-level")
	cmd.Flags().BoolVarP(&getCmd.clusterLevel, "cluster", "c", false, "List templates defined at the cluster-level")
	cmd.Flags().StringVarP(&getCmd.serviceType, "type", "t", "", "Filter the templates by a service type")

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

	instts, err := c.App().GetInstanceTemplates(c.ns, c.clusterLevel, c.brokerLevel, c.serviceType)
	if err != nil {
		return err
	}

	svcattoutput.WriteInstanceTemplateList(c.Output, instts...)
	return nil
}

func (c *getCmd) get() error {
	instt, err := c.App().GetInstanceTemplate(c.ns, c.clusterLevel, c.brokerLevel, c.name)
	if err != nil {
		return err
	}

	svcattoutput.WriteInstanceTemplateList(c.Output, instt)
	return nil
}
