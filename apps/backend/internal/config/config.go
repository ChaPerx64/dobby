package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	OIDCAuthority           string
	OIDCBackendClientID     string
	OIDCBackendClientSecret string
	BackendPort             string
	AllowedOrigins          []string
}

func Load() Config {
	return Config{
		OIDCAuthority:           requireEnv("OIDC_AUTHORITY"),
		OIDCBackendClientID:     requireEnv("OIDC_BACKEND_CLIENT_ID"),
		OIDCBackendClientSecret: requireEnv("OIDC_BACKEND_CLIENT_SECRET"),
		BackendPort:             requireEnv("BACKEND_PORT"),
		AllowedOrigins:          getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
	}
}

func requireEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s is required but not set", key)
	}
	return value
}

func getEnvAsSlice(key string, fallback []string) []string {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parts := strings.Split(valStr, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
}
