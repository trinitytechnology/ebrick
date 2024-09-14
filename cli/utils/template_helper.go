package utils

import (
	"fmt"
	"text/template"
)

func GenerateFileFromTemplate(filePath string, data any, tempContent string) {
	file, err := CreateFile(filePath)
	if err != nil {
		fmt.Println("Error creating file "+filePath, err)
		return
	}
	defer file.Close()

	// Parse the template file
	tmpl, err := template.New("ebrick").Parse(tempContent)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	tmpl.Execute(file, data)
}
