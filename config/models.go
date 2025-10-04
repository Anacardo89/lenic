package config

import "time"

type Config struct {
	Server  Server  `yaml:"server"`
	Session Session `yaml:"pagination"`
	Token   Token   `yaml:"token"`
	DB      DB      `yaml:"db"`
	Log     Log     `yaml:"logging"`
	Img     Img     `yaml:"img"`
	Mail    Mail
}

type Server struct {
	Host            string        `env:"HOST" envDefault:"localhost"`
	Port            string        `env:"PORT" envDefault:"8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`     // seconds
	WriteTimeout    time.Duration `yaml:"write_timeout"`    // seconds
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"` // seconds
}

type Session struct {
	Secret   string        `env:"SESSION_SECRET" envDefault:"session-secret"`
	Duration time.Duration `yaml:"duration"` // hours
}

type Token struct {
	Secret   string        `env:"TOKEN_SECRET" envDefault:"token-secret"`
	Duration time.Duration `yaml:"duration"` // minutes
}

type DB struct {
	DSN             string        `env:"DB_DSN" envDefault:"postgres://user:pass@db:5432/dbname?sslmode=disable"`
	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`  // minutes
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"` // minutes
}

type Log struct {
	Path       string `env:"LOG_PATH" envDefault:"/lenic/logs"`
	File       string `env:"LOG_FILE" envDefault:"lenic.log"`
	Level      string `env:"LOG_LEVEL" envDefault:"info"`
	MaxSize    int    `yaml:"max_size"` // MB
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"` // days
	Compress   bool   `yaml:"compress"`
}

type Img struct {
	BasePath      string `env:"IMG_PATH" envDefault:"/lenic/img"`
	OriginalsDir  string `yaml:"originals_dir"`
	PreviewsDir   string `yaml:"previews_dir"`
	PreviewWidth  int    `yaml:"preview_width"`  // pixels
	PreviewHeight int    `yaml:"preview_height"` // pixels
	JPEGQuality   int    `yaml:"jpeg_quality"`
}

type Mail struct {
	Host string `env:"MAIL_HOST" envDefault:"smtp.gmail.com"`
	Port int    `env:"MAIL_PORT" envDefault:"587"`
	User string `env:"MAIL_USER" envDefault:"user"`
	Pass string `env:"MAIL_PASS" envDefault:"passs"`
}
