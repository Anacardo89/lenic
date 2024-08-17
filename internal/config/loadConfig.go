package config

import (
	_ "embed"

	"gopkg.in/yaml.v3"

	"github.com/Anacardo89/tpsi25_blog/internal/server"
	"github.com/Anacardo89/tpsi25_blog/pkg/auth"
	"github.com/Anacardo89/tpsi25_blog/pkg/db"
	"github.com/Anacardo89/tpsi25_blog/pkg/rabbitmq"
)

//go:embed dbConfig.yaml
var dbYaml []byte

//go:embed serverConfig.yaml
var serverYaml []byte

//go:embed rabbitConfig.yaml
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

func LoadServerConfig() (*server.Config, error) {
	var config server.Config
	err := yaml.Unmarshal(serverYaml, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadRabbitConfig() (*rabbitmq.Config, error) {
	var config rabbitmq.Config
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
