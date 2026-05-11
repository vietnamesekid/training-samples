package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AppConfig chứa tất cả cấu hình của ứng dụng
type AppConfig struct {
	// Server
	Host string
	Port int

	// Database
	DatabaseURL    string
	MaxConnections int

	// Timeouts
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// Feature flags
	EnableMetrics bool
	LogLevel      string
}

// getEnv đọc env var với fallback default
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}

// LoadConfig đọc config từ environment variables
// NGUYÊN TẮC: 12-factor app — config từ environment, không hardcode
func LoadConfig() AppConfig {
	return AppConfig{
		Host:           getEnv("APP_HOST", "0.0.0.0"),
		Port:           getEnvInt("APP_PORT", 8080),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://localhost/mydb"),
		MaxConnections: getEnvInt("DB_MAX_CONNS", 25),
		ReadTimeout:    getEnvDuration("READ_TIMEOUT", 30*time.Second),
		WriteTimeout:   getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
		EnableMetrics:  getEnvBool("ENABLE_METRICS", true),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func (c AppConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func demoConfig() {
	// Set một vài env vars để minh họa
	os.Setenv("APP_PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENABLE_METRICS", "false")
	defer func() {
		os.Unsetenv("APP_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENABLE_METRICS")
	}()

	cfg := LoadConfig()

	fmt.Printf("  Host: %s\n", cfg.Host)
	fmt.Printf("  Port: %d (from APP_PORT env)\n", cfg.Port)
	fmt.Printf("  LogLevel: %s (from LOG_LEVEL env)\n", cfg.LogLevel)
	fmt.Printf("  EnableMetrics: %v (from ENABLE_METRICS env)\n", cfg.EnableMetrics)
	fmt.Printf("  ReadTimeout: %v (default)\n", cfg.ReadTimeout)
	fmt.Printf("  MaxConnections: %d (default)\n", cfg.MaxConnections)

	if err := cfg.Validate(); err != nil {
		fmt.Printf("  Config invalid: %v\n", err)
	} else {
		fmt.Println("  Config valid")
	}

	// 12-Factor App principles:
	fmt.Println("\n  12-Factor Config principles:")
	fmt.Println("  - Store config in environment variables")
	fmt.Println("  - Never hardcode credentials or URLs")
	fmt.Println("  - Provide sensible defaults for development")
	fmt.Println("  - Validate config at startup, fail fast")
}
