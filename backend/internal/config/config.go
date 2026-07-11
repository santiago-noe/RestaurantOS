package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Env         string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:santiago09@localhost:5432/restaurantos?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-cambiar-en-produccion"),
		Port:        getEnv("PORT", "8080"),
		Env:         getEnv("ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
