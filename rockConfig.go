package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"sync"

	"go.yaml.in/yaml/v3"
)

type CommandConfig struct {
	Token            string `yaml:"token"`
	Command          string `yaml:"command"`
	UseGhcrLogin     bool   `yaml:"use-ghcr-login"`
	DontOverrideHome bool   `yaml:"dont-override-home"`
}

func (c *CommandConfig) getBashCommand(request ExecuteRequest) (string, error) {
	cmd := c.Command

	if c.UseGhcrLogin {
		if request.githubActor == "" || request.githubToken == "" {
			return "", fmt.Errorf("This command requires both X-Github-Token and X-Github-Action")
		}
		if err := request.validateGithubActor(); err != nil {
			return "", err
		}
		if err := request.validateGithubToken(); err != nil {
			return "", err
		}
		cmd = "echo " + request.githubToken + " | docker login ghcr.io -u " + request.githubActor + " --password-stdin && " + cmd + " && docker logout ghcr.io"
	}

	return cmd, nil
}

type RockConfig struct {
	Commands map[string]CommandConfig `yaml:"commands"`
}

type RockConfigs struct {
	configs map[string]RockConfig
	mutex   sync.RWMutex
}

var pathConfigs = RockConfigs{
	configs: make(map[string]RockConfig),
	mutex:   sync.RWMutex{},
}

func reloadRock(directory string) error {
	pathConfigs.mutex.Lock()
	defer pathConfigs.mutex.Unlock()
	log.Println("Loading rock config:", directory)

	filepath := config.joinRockDir(directory, ".rhyolite.yaml")

	data, err := os.ReadFile(filepath)
	if err != nil {
		delete(pathConfigs.configs, directory)
		log.Println("No rock config found, removing reference to:", directory)
		return nil // We silently fail if the file doesnt exist, this is intended.
	}

	config := RockConfig{}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("YAML unmarshal failed: %w", err)
	}
	for _, cmd := range config.Commands {
		if len(cmd.Token) < 10 {
			log.Printf("Warning: Command [%s] for Rock [%s] has very short token, please generate a proper token\n", cmd.Command, directory)
		}
	}
	pathConfigs.configs[directory] = config
	return nil
}

func findCommandConfig(request ExecuteRequest) (CommandConfig, error) {
	pathConfigs.mutex.RLock()
	defer pathConfigs.mutex.RUnlock()

	pathConfig, ok := pathConfigs.configs[request.rock]
	if !ok {
		return CommandConfig{}, fmt.Errorf("No config found with name")
	}

	cmd, ok := pathConfig.Commands[request.commandName]
	if !ok {
		return CommandConfig{}, fmt.Errorf("No command found with name")
	}

	requestToken := []byte(request.token)
	configToken := []byte(cmd.Token)

	if subtle.ConstantTimeCompare(requestToken, configToken) != 1 {
		return CommandConfig{}, fmt.Errorf("Invalid token for command")
	}
	return cmd, nil

}
