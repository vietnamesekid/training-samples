// Package config — internal package, only usable within this module
// CANNOT be imported from another module (Go enforces this at compile time)
package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	Host     string
	Port     int
	DBHost   string
	DBPort   int
	DBName   string
	LogLevel string
}

// DatabaseURL returns the connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%d/%s", c.DBHost, c.DBPort, c.DBName)
}

// Load reads config from environment variables with fallback defaults
func Load() *Config {
	return &Config{
		Host:     getEnv("APP_HOST", "localhost"),
		Port:     8080,
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   5432,
		DBName:   getEnv("DB_NAME", "myapp"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
