package main

import (
	"github.com/linkifysoft/ebrick"
	"github.com/linkifysoft/ebrick/examples/modules/environment"
	"github.com/linkifysoft/ebrick/examples/modules/tenant"
)

func main() {
	app := ebrick.NewApplication()
	app.RegisterModules(&tenant.TenantModule{}, &environment.EnvironmentModule{})
	app.Start()
}
