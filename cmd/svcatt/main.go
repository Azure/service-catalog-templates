// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package main

import (
	"fmt"
	"os"

	"github.com/Azure/service-catalog-templates/cmd/svcatt/binding-template"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/command"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/instance-template"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/templated-binding"
	"github.com/Azure/service-catalog-templates/cmd/svcatt/templated-instance"
	"github.com/Azure/service-catalog-templates/pkg/svcatt"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/binding"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/broker"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/class"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/instance"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/plan"
	"github.com/kubernetes-incubator/service-catalog/cmd/svcat/plugin"
	"github.com/kubernetes-incubator/service-catalog/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	cmd := buildRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func buildRootCommand() *cobra.Command {
	// root command context
	cxt := svcattcommand.NewContext()

	// root command flags
	var opts struct {
		Version     bool
		KubeConfig  string
		KubeContext string
	}

	cmd := &cobra.Command{
		Use:          "svcat",
		Short:        "The Kubernetes Service Catalog Command-Line Interface (CLI) with support for templates",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Enable tests to swap the output
			cxt.Output = cmd.OutOrStdout()

			// Initialize flags from environment variables
			bindViperToCobra(cxt.Viper, cmd)

			app, err := svcatt.NewApp(opts.KubeConfig, opts.KubeContext)
			cxt.SetApp(app)

			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				printVersion(cxt)
				return nil
			}

			fmt.Fprint(cxt.Output, cmd.UsageString())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&opts.Version, "version", "v", false, "Show the application version")

	if plugin.IsPlugin() {
		plugin.BindEnvironmentVariables(cxt.Viper)
	} else {
		cmd.PersistentFlags().StringVar(&opts.KubeContext, "kube-context", "", "name of the kube context to use")
		cmd.PersistentFlags().StringVar(&opts.KubeConfig, "kubeconfig", "", "path to kubeconfig file. Overrides $KUBECONFIG")
	}

	cmd.AddCommand(newGetCmd(cxt))
	cmd.AddCommand(newDescribeCmd(cxt))
	cmd.AddCommand(templatedinstance.NewProvisionCmd(cxt))
	cmd.AddCommand(templatedinstance.NewDeprovisionCmd(cxt))
	cmd.AddCommand(templatedbinding.NewBindCmd(cxt))
	cmd.AddCommand(templatedbinding.NewUnbindCmd(cxt))
	cmd.AddCommand(newSyncCmd(cxt))
	cmd.AddCommand(newInstallCmd(cxt))
	cmd.AddCommand(newTouchCmd(cxt))

	return cmd
}

func printVersion(cxt *svcattcommand.Context) {
	fmt.Fprintf(cxt.Output, "svcatt %s (with Templates support)\n", pkg.VERSION)
}

func newSyncCmd(cxt *svcattcommand.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Short:   "Syncs service catalog for a service broker",
		Aliases: []string{"relist"},
	}
	cmd.AddCommand(broker.NewSyncCmd(cxt.Context))

	return cmd
}

func newGetCmd(cxt *svcattcommand.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "List a resource, optionally filtered by name",
	}
	cmd.AddCommand(binding.NewGetCmd(cxt.Context))
	cmd.AddCommand(broker.NewGetCmd(cxt.Context))
	cmd.AddCommand(class.NewGetCmd(cxt.Context))
	cmd.AddCommand(instance.NewGetCmd(cxt.Context))
	cmd.AddCommand(plan.NewGetCmd(cxt.Context))

	cmd.AddCommand(templatedinstance.NewGetCmd(cxt))
	cmd.AddCommand(templatedbinding.NewGetCmd(cxt))
	cmd.AddCommand(instancetemplate.NewGetCmd(cxt))
	cmd.AddCommand(bindingtemplate.NewGetCmd(cxt))

	return cmd
}

func newDescribeCmd(cxt *svcattcommand.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a specific resource",
	}
	cmd.AddCommand(binding.NewDescribeCmd(cxt.Context))
	cmd.AddCommand(broker.NewDescribeCmd(cxt.Context))
	cmd.AddCommand(class.NewDescribeCmd(cxt.Context))
	cmd.AddCommand(instance.NewDescribeCmd(cxt.Context))
	cmd.AddCommand(plan.NewDescribeCmd(cxt.Context))

	cmd.AddCommand(templatedinstance.NewDescribeCmd(cxt))
	cmd.AddCommand(templatedbinding.NewDescribeCmd(cxt))
	cmd.AddCommand(instancetemplate.NewDescribeCmd(cxt))
	cmd.AddCommand(bindingtemplate.NewDescribeCmd(cxt))

	return cmd
}

func newInstallCmd(cxt *svcattcommand.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "install",
	}
	cmd.AddCommand(plugin.NewInstallCmd(cxt.Context))

	return cmd
}

func newTouchCmd(cxt *svcattcommand.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "touch",
		Short: "Update a resource to trigger reprocessing",
	}
	cmd.AddCommand(instance.NewTouchCommand(cxt.Context))
	return cmd
}

// Bind the viper configuration back to the cobra command flags.
// Allows us to interact with the cobra flags normally, and while still
// using viper's automatic environment variable binding.
func bindViperToCobra(vip *viper.Viper, cmd *cobra.Command) {
	vip.BindPFlags(cmd.Flags())
	vip.AutomaticEnv()
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && vip.IsSet(f.Name) {
			cmd.Flags().Set(f.Name, vip.GetString(f.Name))
		}
	})
}
