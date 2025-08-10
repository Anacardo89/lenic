package server

import (
	"database/sql"
	"net/http"
)

type Config struct {
	Host      string `yaml:"host"`
	ProxyPORT string `yaml:"proxyPort"`
	HttpPORT  string `yaml:"httpPort"`
	HttpsPORT string `yaml:"httpsPort"`
}

type Server struct {
	http.Server
	DB *sql.DB
	SessionStore
	// other dependencies
}

var (
	Server *Config
)
