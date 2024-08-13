package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"

	"github.com/Anacardo89/tpsi25_blog/auth"
	"github.com/Anacardo89/tpsi25_blog/internal/rabbit"
	"github.com/Anacardo89/tpsi25_blog/internal/routes"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
)

//go:embed dbConfig.yaml
var dbYaml []byte

//go:embed serverConfig.yaml
var serverYaml []byte

//go:embed serverConfig.yaml
var rabbitYaml []byte

//go:embed sessionConfig.yaml
var sessionConfig []byte

func LoadDBConfig() (*db.Config, error) {
	var config db.Config
	err := yaml.Unmarshal(dbYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadServerConfig() (*routes.Config, error) {
	var config routes.Config
	err := yaml.Unmarshal(serverYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadRabbitConfig() (*rabbit.Config, error) {
	var config rabbit.Config
	err := yaml.Unmarshal(rabbitYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadSessionConfig() (*auth.Config, error) {
	var config *auth.Config
	err := yaml.Unmarshal(sessionConfig, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
