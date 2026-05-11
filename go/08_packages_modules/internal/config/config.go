// Package config — internal package, chỉ dùng trong module này
// KHÔNG thể import từ module khác (Go enforce ở compile time)
package config

import (
	"fmt"
	"os"
)

// Config chứa cấu hình ứng dụng
type Config struct {
	Host     string
	Port     int
	DBHost   string
	DBPort   int
	DBName   string
	LogLevel string
}

// DatabaseURL trả về connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%d/%s", c.DBHost, c.DBPort, c.DBName)
}

// Load đọc config từ environment variables với fallback defaults
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
