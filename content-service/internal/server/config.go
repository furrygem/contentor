package server

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ListenAddr string `json:"listen_addr" yaml:"listen_addr" env:"CONTENT_SERVICE_SERVER_LISTEN_ADDR"`
}

func (c *Config) FromENV() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}

func NewConfig() *Config {
	return &Config{
		ListenAddr: "127.0.0.1:8080",
	}
}
