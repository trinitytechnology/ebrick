package main

import (
	"github.com/trinitytechnology/ebrick"
	"github.com/trinitytechnology/ebrick/examples/modules/environment"
)

func main() {
	app := ebrick.NewApplication()
	app.RegisterModules(&environment.EnvironmentModule{})
	app.Start()
}
