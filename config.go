package main

import (
	"log"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	WebAddress string `yaml:"web_address"`
	RockDir    string `yaml:"rock_dir"`
}

func (c *Config) joinRockDir(paths ...string) string {
	o := c.RockDir
	for _, path := range paths {
		o = filepath.Join(o, path)
	}

	return o
}

var config Config

func LoadConfig() {
	address := "/etc/rhyolite/config.yaml"
	configFile, err := os.ReadFile(address)
	if err != nil {
		log.Fatalf("Unable to read config:", err.Error())
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse config:", err.Error())
	}
}
