// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package templatedinstance

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/output"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/parameters"
	"github.com/spf13/cobra"
)

type provisonCmd struct {
	*svcattcommand.Context

	ns           string
	instanceName string
	serviceType  string
	className    string
	planName     string
	rawParams    []string
	jsonParams   string
	params       interface{}
	rawSecrets   []string
	secrets      map[string]string
}

// NewProvisionCmd builds a "svcat provision" command
func NewProvisionCmd(cxt *svcattcommand.Context) *cobra.Command {
	provisionCmd := &provisonCmd{Context: cxt}
	cmd := &cobra.Command{
		Use:   "provision NAME --type SERVICE_TYPE",
		Short: "Create a new instance of a service type",
		Example: `
  svcat provision mysql-instance --type mysqldb
  svcat provision mysql-instance --type mysqldb -p location=eastus
  svcat provision mysql-instance --type mysqldb -s mysecret[dbparams]
  svcat provision secure-mysql-instance --type mysqldb --params-json '{
    "encrypt" : true,
    "firewallRules" : [
        {
            "name": "AllowSome",
            "startIPAddress": "75.70.113.50",
            "endIPAddress" : "75.70.113.131"
        },
        {
            "name": "AllowMore",
            "startIPAddress": "13.54.0.0",
            "endIPAddress" : "13.56.0.0"
        }
    ]
  }
  svcat provision wordpress-mysql-instance --class mysqldb --plan free
'
`,
		PreRunE: command.PreRunE(provisionCmd),
		RunE:    command.RunE(provisionCmd),
	}
	cmd.Flags().StringVarP(&provisionCmd.ns, "namespace", "n", "",
		"The namespace in which to create the resource")
	cmd.Flags().StringVar(&provisionCmd.serviceType, "type", "",
		"The service type. Either --type or --class and --plan is required")
	cmd.Flags().StringVar(&provisionCmd.className, "class", "",
		"The class name. Either --type or --class and --plan is required")
	cmd.Flags().StringVar(&provisionCmd.planName, "plan", "",
		"The plan name. Either --type or --class and --plan is required")
	cmd.Flags().StringSliceVarP(&provisionCmd.rawParams, "param", "p", nil,
		"Additional parameter to use when provisioning the service, format: NAME=VALUE. Cannot be combined with --params-json")
	cmd.Flags().StringSliceVarP(&provisionCmd.rawSecrets, "secret", "s", nil,
		"Additional parameter, whose value is stored in a secret, to use when provisioning the service, format: SECRET[KEY]")
	cmd.Flags().StringVar(&provisionCmd.jsonParams, "params-json", "",
		"Additional parameters to use when provisioning the service, provided as a JSON object. Cannot be combined with --param")
	return cmd
}

func (c *provisonCmd) Validate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("an instance name is required")
	}
	c.instanceName = args[0]

	if c.ns == "" {
		c.ns = c.App().CurrentNamespace
	}

	if c.serviceType == "" {
		if c.planName == "" || c.className == "" {
			return fmt.Errorf("either --type or --class and --plan is required")
		}
	} else {
		if c.planName != "" || c.className != "" {
			return fmt.Errorf("--type cannot be used with --class and --plan")
		}
	}

	var err error

	if c.jsonParams != "" && len(c.rawParams) > 0 {
		return fmt.Errorf("--params-json cannot be used with --param")
	}

	if c.jsonParams != "" {
		c.params, err = parameters.ParseVariableJSON(c.jsonParams)
		if err != nil {
			return fmt.Errorf("invalid --params value (%s)", err)
		}
	} else {
		c.params, err = parameters.ParseVariableAssignments(c.rawParams)
		if err != nil {
			return fmt.Errorf("invalid --param value (%s)", err)
		}
	}

	c.secrets, err = parameters.ParseKeyMaps(c.rawSecrets)
	if err != nil {
		return fmt.Errorf("invalid --secret value (%s)", err)
	}

	return nil
}

func (c *provisonCmd) Run() error {
	return c.Provision()
}

func (c *provisonCmd) Provision() error {
	tinst, err := c.App().Provision(c.ns, c.instanceName, c.serviceType, c.className, c.planName, c.params, c.secrets)
	if err != nil {
		return err
	}

	svcattoutput.WriteTemplatedInstanceDetails(c.Output, tinst)

	return nil
}
