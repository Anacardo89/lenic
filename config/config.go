package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var configYaml []byte

func LoadConfig() (*Config, error) {
	var config Config
	err := yaml.Unmarshal(configYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
