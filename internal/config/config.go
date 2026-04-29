package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port         string
	GRPCPort     string
	ReadTimeout  int
	WriteTimeout int
}

func LoadServerConfig() *ServerConfig {
	_ = godotenv.Load()

	return &ServerConfig{
		Port:         getEnv("CONFIGLINTER_PORT", "8080"),
		GRPCPort:     getEnv("CONFIGLINTER_GRPC_PORT", "9090"),
		ReadTimeout:  getEnvInt("CONFIGLINTER_READ_TIMEOUT", 10),
		WriteTimeout: getEnvInt("CONFIGLINTER_WRITE_TIMEOUT", 10),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
