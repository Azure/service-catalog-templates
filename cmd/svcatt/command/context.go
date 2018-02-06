// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package svcattcommand

import (
	"github.com/Azure/service-catalog-templates/pkg/svcatt"
	svcatcommand "github.com/kubernetes-incubator/service-catalog/cmd/svcat/command"
	"github.com/spf13/viper"
)

// Context is ambient data necessary to run any svcatt command.
type Context struct {
	*svcatcommand.Context

	// svcat application, the library behind this cli
	app *svcatt.App
}

// NewContext for the svcatt cli
func NewContext() *Context {
	return &Context{
		Context: &svcatcommand.Context{
			Viper: viper.New(),
		},
	}
}

func (cxt *Context) SetApp(app *svcatt.App) {
	cxt.app = app
	cxt.Context.App = app.ServiceCatalogApp
}

func (cxt *Context) App() *svcatt.App {
	return cxt.app
}
