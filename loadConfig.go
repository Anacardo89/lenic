package main

import (
	_ "embed"

	"gopkg.in/yaml.v3"

	"github.com/Anacardo89/tpsi25_blog.git/auth"
	"github.com/Anacardo89/tpsi25_blog.git/db"
)

//go:embed dbConfig.yaml
var dbYaml []byte

//go:embed serverConfig.yaml
var serverYaml []byte

//go:embed serverConfig.yaml
var rabbitYaml []byte

//go:embed sessionConfig.yaml
var sessionConfig []byte

type RabbitConfig struct {
	MQHost string `yaml:"mqHost"`
	MQPort string `yaml:"mqPort"`
}

type ServerConfig struct {
	ProxyPORT string `yaml:"proxyPort"`
	HttpPORT  string `yaml:"httpPort"`
	HttpsPORT string `yaml:"httpsPort"`
}

func loadDBConfig() (*db.DBConfig, error) {
	var dbConfig db.DBConfig
	err := yaml.Unmarshal(dbYaml, &dbConfig)
	if err != nil {
		return nil, err
	}
	return &dbConfig, nil
}

func loadServerConfig() (*ServerConfig, error) {
	var serverConfig ServerConfig
	err := yaml.Unmarshal(serverYaml, &serverConfig)
	if err != nil {
		return nil, err
	}
	return &serverConfig, nil
}

func loadRabbitConfig() (*RabbitConfig, error) {
	var rabbitConfig RabbitConfig
	err := yaml.Unmarshal(rabbitYaml, &rabbitConfig)
	if err != nil {
		return nil, err
	}
	return &rabbitConfig, nil
}

func loadSessionConfig() (*auth.SessionConfig, error) {
	var sessConfig *auth.SessionConfig
	err := yaml.Unmarshal(sessionConfig, &sessConfig)
	if err != nil {
		return nil, err
	}
	return sessConfig, nil
}
