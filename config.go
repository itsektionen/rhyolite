package main

import (
	"crypto/subtle"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.yaml.in/yaml/v3"
)

type CommandConfig struct {
	Token   string `yaml:"token"`
	Command string `yaml:"command"`
}

type PathConfig struct {
	Commands []CommandConfig `yaml:"commands"`
}

type PathConfigs struct {
	configs map[string]PathConfig
	mutex   sync.RWMutex
}

var pathConfigs = PathConfigs{
	configs: make(map[string]PathConfig),
	mutex:   sync.RWMutex{},
}

func reloadPath(directory string) error {
	pathConfigs.mutex.Lock()
	defer pathConfigs.mutex.Unlock()
	println("Loading rock:", directory)

	filepath := filepath.Join(RhyoliteDir, directory, ".rhyolite.yaml")

	data, err := os.ReadFile(filepath)
	if err != nil {
		delete(pathConfigs.configs, directory)
		return nil // We silently fail if the file doesnt exist, this is intended.
	}

	config := PathConfig{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("YAML unmarshal failed: %w", err)
	}
	pathConfigs.configs[directory] = config
	return nil
}

func findCommandConfig(request ExecuteRequest) (CommandConfig, error) {
	pathConfigs.mutex.RLock()
	defer pathConfigs.mutex.RUnlock()

	pathConfig, ok := pathConfigs.configs[request.directory]
	if !ok {
		return CommandConfig{}, fmt.Errorf("No config found with name")
	}

	for _, cmd := range pathConfig.Commands {
		requestToken := []byte(request.token)
		configToken := []byte(cmd.Token)

		if subtle.ConstantTimeCompare(requestToken, configToken) == 1 {
			return cmd, nil
		}
	}

	return CommandConfig{}, fmt.Errorf("No command matching the provided token was found")
}
