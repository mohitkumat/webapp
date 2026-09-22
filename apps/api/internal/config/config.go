package config

import "os"

type Config struct {
	Host           string
	Port           string
	DatabaseURL    string
	AllowedOrigins string
}

func Load() Config {
	return Config{
		Host:           getEnv("HOST", "127.0.0.1"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"),
		AllowedOrigins: getEnv("ALLOWED_ORIGIN", "http://127.0.0.1:3000,http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
