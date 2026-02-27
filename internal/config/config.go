package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppEnv      string
	DatabaseURL string
	Port        int
	LogLevel    logrus.Level
}

func LoadConfig() (*Config, error) {
	logLevel, err := getEnvAsLogLevel("LOG_LEVEL", logrus.InfoLevel)
	if err != nil {
		return nil, err
	}

	port, err := getEnvAsInt("PORT", 8080)
	if err != nil {
		return nil, err
	}

	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Port:        port,
		LogLevel:    *logLevel,
	}, nil
}

func LoadEnv() {
	godotenv.Load(".env." + getEnv("APP_ENV", "development"))
	godotenv.Load(".env")
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}

		return intValue, nil
	}

	return defaultValue, nil
}

func getEnvAsLogLevel(key string, defaultValue logrus.Level) (*logrus.Level, error) {
	if value, exists := os.LookupEnv(key); exists {
		level, err := logrus.ParseLevel(value)
		if err != nil {
			return nil, err
		}

		return &level, nil
	}

	return &defaultValue, nil
}
