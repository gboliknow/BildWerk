package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTSecret  string
}

var Envs = InitializeConfig()

func InitializeConfig() Config {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println("Error:", err)
	}

	return Config{
		Port:       getEnv("PORT", "8082"),
		DBUser:     getEnv("DB_USER", "user3"),
		DBPassword: getEnv("DB_PASSWORD", "pass3"),
		DBName:     getEnv("DB_NAME", "bildwerk"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "8082")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
