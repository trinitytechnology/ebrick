package main

import (
	"github.com/linkifysoft/ebrick"
	"github.com/linkifysoft/ebrick/examples/modules/environment"
)

func main() {
	app := ebrick.NewApplication()
	app.RegisterModules(&environment.EnvironmentModule{})
	app.Start()
}
