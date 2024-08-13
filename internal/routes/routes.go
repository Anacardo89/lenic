package routes

type Config struct {
	ProxyPORT string `yaml:"proxyPort"`
	HttpPORT  string `yaml:"httpPort"`
	HttpsPORT string `yaml:"httpsPort"`
}

var (
	Server *Config
)
