package minioHandler

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Workers         int    `env:"CONTENT_SERVICE_MINIO_WORKERS"`
	ChannelCapacity int    `env:"CONTENT_SERVICE_MINIO_CHANNEL_CAPACITY"`
	BucketName      string `env:"CONTENT_SERVICE_MINIO_BUCKET_NAME"`
	Endpoint        string `env:"CONTENT_SERVICE_MINIO_ENDPOINT"`
	AccessKeyID     string `env:"CONTENT_SERVICE_MINIO_ACCESS_KEY"`
	SecretKeyID     string `env:"CONTENT_SERVICE_MINIO_SECRET_KEY"`
	UseSSL          bool   `env:"CONTENT_SERVICE_MINIO_USE_SSL"`
}

func (c *Config) FromENV() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	return nil
}

func NewConfig() *Config {
	return &Config{
		BucketName:      "files",
		Workers:         1,
		ChannelCapacity: 1,
		UseSSL:          false,
		Endpoint:        "localhost:9000",
	}
}
