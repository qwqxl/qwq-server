package config

import (
	"qwqserver/pkg/defaultvalue"
)

func NewListen() (*Listen, error) {
	c, err := New()
	return c.Listen, err
}

func NewDatabase() (*Database, error) {
	c, err := New()
	return c.Database, err
}

func NewCache() (*Cache, error) {
	c, err := New()
	return c.Cache, err
}

func NewAdminUser() (*AdminUser, error) {
	c, err := New()

	return c.AdminUser, err
}

func NewServer() (*Server, error) {
	c, err := New()
	return c.Server, err
}

func NewAuth() (*Auth, error) {
	c, err := New()
	return c.Auth, err
}

func Default() (*Config, error) {
	defaultConfig := &Config{}
	err := defaultvalue.SetDefaults(defaultConfig)
	if err != nil {
		return nil, err
	}
	return defaultConfig, nil
}
