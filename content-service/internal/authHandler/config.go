package authHandler

import (
	"github.com/caarlos0/env/v6"
)

type AuthHandlerConfig struct {
	EdDSAPublicKeyFile string `json:"eddsa_public_key_file" env:"CONTENT_SERVICE_AUTH_EDDSA_PUBLIC_KEY_FILE"`
}

func NewConfig() *AuthHandlerConfig {
	return &AuthHandlerConfig{}
}

func (c *AuthHandlerConfig) FromEnv() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}
