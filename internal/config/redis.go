package config

import "time"

type Redis struct {
	Addr         string        `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379" qwq-default:"localhost:6379"`
	Password     string        `yaml:"password" env:"REDIS_PASSWORD" env-default:"" qwq-default:"123456"`
	DB           int           `yaml:"db" env:"REDIS_DB" env-default:"0" qwq-default:"0"`
	PoolSize     int           `yaml:"pool_size" env:"REDIS_POOL_SIZE" env-default:"10" qwq-default:"10"`
	MinIdleConns int           `yaml:"min_idle_conns" env:"REDIS_MIN_IDLE_CONNS" env-default:"5" qwq-default:"5"`
	MaxRetries   int           `yaml:"max_retries" env:"REDIS_MAX_RETRIES" env-default:"3" qwq-default:"3"`
	DialTimeout  time.Duration `yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT" env-default:"5s" qwq-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"REDIS_READ_TIMEOUT" env-default:"3s" qwq-default:"3s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"REDIS_WRITE_TIMEOUT" env-default:"3s" qwq-default:"3s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"REDIS_IDLE_TIMEOUT" env-default:"5m" qwq-default:"5m"`
}
