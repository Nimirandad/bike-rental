package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SQLitePath string
	Port       string
	LogLevel   string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		SQLitePath: getEnvDefault("SQLITE_PATH", SQLitePath),
		Port:       getEnvDefault("HTTP_PORT", HTTPPort),
		LogLevel:   getEnvDefault("LOG_LEVEL", LogLevel),
	}
}

func getEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
