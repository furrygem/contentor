package server

import "github.com/furrygem/contentor/content-service/pkg/utils"

type Config struct {
	ListenAddr string `json:"listen_addr" yaml:"listen_addr"`
}

func NewConfig() *Config {
	return &Config{
		ListenAddr: utils.GetEnvDefault("ListenAddr", "127.0.0.1", false),
	}
}
