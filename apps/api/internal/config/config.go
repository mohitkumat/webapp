package config

import "os"

type Config struct {
	Host          string
	Port          string
	DatabaseURL   string
	AllowedOrigin string
}

func Load() Config {
	return Config{
		// Default 127.0.0.1 for bare-metal `go run` (loopback-only, per the local
		// workstation security requirement). Inside a container this must be
		// overridden to 0.0.0.0 — Docker's port publishing proxies to the
		// container's network interface, not its internal loopback, so a
		// container process bound to 127.0.0.1 is unreachable from the host even
		// with a published port. The host-side restriction to 127.0.0.1 is
		// instead enforced by compose's "127.0.0.1:PORT:PORT" port mapping.
		Host:          getEnv("HOST", "127.0.0.1"),
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@127.0.0.1:5432/appdb?sslmode=disable"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://127.0.0.1:3000,http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
