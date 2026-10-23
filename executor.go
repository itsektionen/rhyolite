package main

import (
	"crypto/subtle"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Token   string `yaml:"token"`
	Command string `yaml:"command"`
}

func ReadConfig(request ExecuteRequest) (Config, error) {

	config := Config{}

	filepath := filepath.Join(RhyoliteDir, request.directory, "rhyolite.yaml")

	data, err := os.ReadFile(filepath)
	if err != nil {
		return config, fmt.Errorf("Unable to read", filepath, err.Error())
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("Yaml unmarshal failed Data \"", data, "\" error:", err.Error())
	}

	requestToken := []byte(request.token)
	configToken := []byte(config.Token)

	if subtle.ConstantTimeCompare(requestToken, configToken) != 1 {
		return config, fmt.Errorf("Invalid token")
	}

	return config, nil
}

func Execute(config Config, request ExecuteRequest) ([]byte, error) {
	cmd := exec.Command("bash", "-c", config.Command)
	cmd.Dir = filepath.Join(RhyoliteDir, request.directory)
	return cmd.Output()
}
