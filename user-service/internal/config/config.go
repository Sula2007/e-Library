package config

import (
	"os"
)

type Config struct {
	GRPCPort   string
	DBUrl      string
	RedisAddr  string
	JWTSecret  string
	NATSUrl    string
	SMTPHost   string
	SMTPPort   string
	SMTPUser   string
	SMTPPass   string
}

func Load() *Config {
	return &Config{
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		DBUrl:     getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		NATSUrl:   getEnv("NATS_URL", "nats://localhost:4222"),
		SMTPHost:  getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:  getEnv("SMTP_PORT", "587"),
		SMTPUser:  getEnv("SMTP_USER", ""),
		SMTPPass:  getEnv("SMTP_PASS", ""),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}