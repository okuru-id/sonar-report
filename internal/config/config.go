package config

import (
	"os"
	"strconv"
)

type Config struct {
	// SonarQube
	SonarQubeURL   string
	SonarQubeToken string

	// Server
	ServerPort string

	// Admin Auth
	AdminUsername string
	AdminPassword string

	// Session
	SessionSecret string

	// Report Storage
	ReportStoragePath   string
	ReportRetentionDays int
}

func Load() *Config {
	return &Config{
		SonarQubeURL:        getEnv("SONARQUBE_URL", "http://localhost:9000"),
		SonarQubeToken:      getEnv("SONARQUBE_TOKEN", ""),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		AdminUsername:       getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "admin"),
		SessionSecret:       getEnv("SESSION_SECRET", "default-secret-key-change-in-production"),
		ReportStoragePath:   getEnv("REPORT_STORAGE_PATH", "./reports"),
		ReportRetentionDays: getEnvInt("REPORT_RETENTION_DAYS", 30),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
