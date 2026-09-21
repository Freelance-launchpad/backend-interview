package confiture

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents a configuration object.
type Config interface {
	Secrets() error
}

func findFilePath(filePath ...string) (string, error) {
	if len(filePath) > 0 {
		return filePath[0], nil
	}

	if len(os.Args[1:]) > 0 {
		return os.Args[1], nil
	}

	return "", fmt.Errorf("no configuration file path specified")
}

// Load loads a conf and call Secrets function if it exists.
func Load(conf any, filePath ...string) error {
	configFilePath, err := findFilePath(filePath...)
	if err != nil {
		return err
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(content, conf); err != nil {
		return err
	}

	switch v := conf.(type) {
	case Config:
		if err := v.Secrets(); err != nil {
			return err
		}
	}

	return nil
}

// LoadSecretFromFile load a secret from a secret file.
func LoadSecretFromFile(envVariableName string) (string, error) {
	path := os.Getenv(envVariableName)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
