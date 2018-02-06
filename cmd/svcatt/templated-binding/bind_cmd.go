// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedbinding

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/parameters"
	"github.com/spf13/cobra"
)

type bindCmd struct {
	*svcattcommand.Context
	ns           string
	instanceName string
	bindingName  string
	secretName   string
	rawParams    []string
	params       map[string]string
	rawSecrets   []string
	secrets      map[string]string
}

// NewBindCmd builds a "svcat bind" command
func NewBindCmd(cxt *svcattcommand.Context) *cobra.Command {
	bindCmd := &bindCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "bind INSTANCE_NAME",
		Short: "Binds an instance's metadata to a secret, which can then be used by an application to connect to the instance",
		Example: `
  svcat bind wordpress
  svcat bind wordpress-mysql-instance --name wordpress-mysql-binding --secret-name wordpress-mysql-secret
`,
		PreRunE: command.PreRunE(bindCmd),
		RunE:    command.RunE(bindCmd),
	}
	cmd.Flags().StringVarP(
		&bindCmd.ns,
		"namespace",
		"n",
		"",
		"The resource namespace",
	)
	cmd.Flags().StringVarP(
		&bindCmd.bindingName,
		"name",
		"",
		"",
		"The name of the binding. Defaults to the name of the instance.",
	)
	cmd.Flags().StringVarP(
		&bindCmd.secretName,
		"secret-name",
		"",
		"",
		"The name of the secret. Defaults to the name of the instance.",
	)
	cmd.Flags().StringSliceVarP(&bindCmd.rawParams, "param", "p", nil,
		"Additional parameter to use when binding the instance, format: NAME=VALUE")
	cmd.Flags().StringSliceVarP(&bindCmd.rawSecrets, "secret", "s", nil,
		"Additional parameter, whose value is stored in a secret, to use when binding the instance, format: SECRET[KEY]")

	return cmd
}

func (c *bindCmd) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("instance is required")
	}
	c.instanceName = args[0]

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	var err error
	c.params, err = parameters.ParseVariableAssignments(c.rawParams)
	if err != nil {
		return fmt.Errorf("invalid --param value (%s)", err)
	}

	c.secrets, err = parameters.ParseKeyMaps(c.rawSecrets)
	if err != nil {
		return fmt.Errorf("invalid --secret value (%s)", err)
	}

	return nil
}

func (c *bindCmd) Run() error {
	return c.bind()
}

func (c *bindCmd) bind() error {
	tbnd, err := c.App().Bind(c.ns, c.bindingName, c.instanceName, c.secretName, c.params, c.secrets)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedBindingDetails(c.Output, tbnd)
	return nil
}
