package config

import "os"

type Config struct {
	GRPCPort string
	DBUrl    string
}

func Load() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", "50054"),
		DBUrl:    getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/paymentdb?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}