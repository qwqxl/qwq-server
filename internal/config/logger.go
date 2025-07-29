package config

type Logger struct {
	Level string `yaml:"level" env:"QWQ_LEVEL_HOME" env-default:"debug" qwq-default:"debug"` // redis地址
}
