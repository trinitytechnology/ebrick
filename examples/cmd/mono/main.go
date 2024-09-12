package main

import (
	"github.com/trinitytechnology/ebrick"
	"github.com/trinitytechnology/ebrick/examples/modules/environment"
	"github.com/trinitytechnology/ebrick/examples/modules/tenant"
)

func main() {
	app := ebrick.NewApplication()
	app.RegisterModules(&tenant.TenantModule{}, &environment.EnvironmentModule{})

	app.Start()
}
