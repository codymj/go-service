package internal

import (
	"context"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxsql"
)

// RootDir is the application's root directory.
var RootDir string

func init() {
	RootDir = fxcore.RootDir(1)
}

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	fxhttpserver.FxHttpServerModule,
	fxsql.FxSQLModule,
	Register(),
	Router(),
)

// Run starts the application, with a provided context.Context.
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}
