package config

import "fmt"

type Listen struct {
	Host          string `yaml:"host" env:"SERVER_HOST" env-default:"localhost" qwq-default:""`
	Port          int    `yaml:"port" env:"SERVER_PORT" env-default:"8080" qwq-default:"5000"`
	LogLevel      string `yaml:"log_level" env:"SERVER_LOG_LEVEL" env-default:"debug" qwq-default:"debug"`
	Mode          string `yaml:"mode" env:"SERVER_MODE" env-default:"debug" default-value:"debug"`
	MaxConcurrent int    `yaml:"max_concurrent" env:"SERVER_MAX_CONCURRENT" env-default:"100" qwq-default:"100"`
}

func (l *Listen) ListenAddress() string {
	if l.Port == 0 {
		l.Port = 5000
	}
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}
