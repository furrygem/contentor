package utils

import (
	"log"
	"os"
)

// GetEnvDefault gets environment variable, if there is no variable in the environment return fallback,
// if required is set to true and the variable isn't set, the execution of the program is stopped
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
