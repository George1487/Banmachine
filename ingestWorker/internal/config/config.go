package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultCountWorkerConcurrency = 4
	defaultPollInterval           = 5 * time.Second
)

type Config struct {
	DBDSN             string
	WorkerConcurrency int
	PollInterval      time.Duration

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func MustConfig() *Config {
	return &Config{
		DBDSN:             mustEnv("DB_DSN"),
		WorkerConcurrency: mustIntEnv("WORKER_CONCURRENCY", defaultCountWorkerConcurrency),
		PollInterval:      mustDurationEnv("POLL_INTERVAL", defaultPollInterval),
		MinioEndpoint:     mustEnv("MINIO_ENDPOINT"),
		MinioAccessKey:    mustEnv("MINIO_ACCESS_KEY"),
		MinioSecretKey:    mustEnv("MINIO_SECRET_KEY"),
		MinioBucket:       mustEnv("MINIO_BUCKET"),
		MinioUseSSL:       mustBoolEnv("MINIO_USE_SSL", false),
	}
}

func mustEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return value
}

func mustIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid %s: %v", key, err))
	}
	if parsed <= 0 {
		panic(fmt.Sprintf("%s must be positive", key))
	}
	return parsed
}

func mustDurationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		panic(fmt.Sprintf("invalid %s: %v", key, err))
	}
	if parsed <= 0 {
		panic(fmt.Sprintf("%s must be positive", key))
	}
	return parsed
}

func mustBoolEnv(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("invalid %s: %v", key, err))
	}
	return parsed
}
