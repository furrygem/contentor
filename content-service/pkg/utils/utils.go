package utils

import (
	"log"
	"os"
)

func GetEnvDefault(key, fallback string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("%s environment variable must be set\n", key)
		os.Exit(1)
	}
	if value == "" {
		return fallback
	}
	return value
}
