package config

import (
	"os"
)

type Config struct {
	Host string
	Port string
}

func New() Config {
	return Config{
		Host: getEnvOrDefault("HOST", "127.0.0.1"),
		Port: getEnvOrDefault("PORT", "8091"),
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
