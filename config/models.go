package config

type Config struct {
	DB     DBConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
	Auth   AuthConfig   `yaml:"auth"`
}

type DBConfig struct {
	DSN  string `yaml:"dsn"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type ServerConfig struct {
	Host      string `yaml:"host"`
	ProxyPort string `yaml:"proxy_port"`
	HTTPPort  string `yaml:"http_port"`
	HTTPSPort string `yaml:"https_port"`
}

type AuthConfig struct {
	SessionPass   string `yaml:"session_pass"`
	SessionExpMin string `yaml:"session_expiration_minutes"`
}
