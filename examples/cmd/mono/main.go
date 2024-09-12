package main

import (
	"github.com/trinitytechnology/ebrick"
	"github.com/trinitytechnology/ebrick/examples/internal/environment"
	"github.com/trinitytechnology/ebrick/examples/internal/tenant"
)

func main() {
	app := ebrick.NewApplication()
	app.RegisterModules(&tenant.TenantModule{}, &environment.EnvironmentModule{})

	app.Start()
}
