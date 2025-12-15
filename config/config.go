package config

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func NewConfig() *Config {
	return &Config{
		Server:  Server{},
		Session: Session{},
		Token:   Token{},
		DB:      DB{},
		Log:     Log{},
		Img: Img{
			ImgDirs:     make(map[string]string),
			PreviewDims: make(map[string]int),
		},
		Mail: Mail{},
	}
}

func LoadConfig() (*Config, error) {
	cfg := NewConfig()
	_ = godotenv.Load()
	appHome := os.Getenv("APP_HOME")
	if appHome == "" {
		return nil, errors.New("APP_HOME not set")
	}
	cfgPath := os.Getenv("CFG_PATH")
	if cfgPath == "" {
		return nil, errors.New("CFG_PATH not set")
	}
	cfgPath = filepath.Join(appHome, cfgPath)
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %s", err)
	}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing env: %s", err)
	}
	return cfg, nil
}
