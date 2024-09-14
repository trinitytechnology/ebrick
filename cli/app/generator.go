package app

import (
	_ "embed"
	"fmt"

	"github.com/trinitytechnology/ebrick/cli/utils"
)

func GenerateApplication(ebrickConfig AppConfig) {

	// Generate the application.yaml file
	GenerateApplicationConfig(ebrickConfig)

	// Generate the main.go file
	GenerateMainFile(ebrickConfig)

	// Generate the docker-compose.yml file
	GenerateDockerComposeFile(ebrickConfig)

	// Generate the go.mod file
	GenerateGoModFile(ebrickConfig)
}

// Embed the template file

//go:embed templates/application.yaml.tmpl
var applicationTemplate string

// GenerateApplicationConfig generates the application.yaml file using a template
func GenerateApplicationConfig(appConfig AppConfig) {

	utils.GenerateFileFromTemplate("application.yaml", appConfig, applicationTemplate)
	fmt.Println("Generated application.yaml successfully.")
}

//go:embed templates/main.go.tmpl
var mainTemplate string

// GenerateMainFile generates the main.go file using a template
func GenerateMainFile(appConfig AppConfig) {
	utils.GenerateFileFromTemplate("cmd/main.go", appConfig, mainTemplate)
}

//go:embed templates/docker-compose.yml.tmpl
var dockerComposeTemplate string

// GenerateMainFile generates the main.go file using a template
func GenerateDockerComposeFile(appConfig AppConfig) {
	utils.GenerateFileFromTemplate("docker-compose.yml", appConfig, dockerComposeTemplate)
}

//go:embed templates/go.mod.tmpl
var goModTemplate string

// GenerateMainFile generates the main.go file using a template
func GenerateGoModFile(appConfig AppConfig) {
	utils.GenerateFileFromTemplate("go.mod", appConfig, goModTemplate)
}
