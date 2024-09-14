package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Helper function to read .ebrick.yaml
func ReadYamlFile[T any](filePath string) (T, error) {
	var ebrickConfig T
	file, err := os.Open(filePath)
	if err != nil {
		return ebrickConfig, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&ebrickConfig)
	return ebrickConfig, err
}

func WriteYamlFile[T any](filePath string, data T) error {
	file, err := CreateFile(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	configData, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling yaml:", err)
		return err
	}

	_, err = file.Write(configData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("Created file:", filePath)
	return nil
}

func CreateFile(filePath string) (*os.File, error) {
	// Create parent directories if they don't exist
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating directories for %s: %w", filePath, err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error creating file %s: %w", filePath, err)
	}

	return file, nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func CreateFolder(folderPath string) error {
	// Create parent directories if they don't exist
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories for %s: %w", folderPath, err)
	}
	return nil
}
