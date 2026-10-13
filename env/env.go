package env

import (
	"errors"
	"strconv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	ErrMissingValue = errors.New("missing value")
	ErrInvalidType  = errors.New("invalid type for variable")
)

type EnvVar interface {
	string | int | bool
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("unable to load environment: %v\n", err)
	}
}

func StringOrDefault(key string, defaultValue string) string {
	value, err := stringFromEnv(key)
	if err != nil {
		return defaultValue
	}
	return value
}

func StringOrFatal(key string) string {
	value, err := stringFromEnv(key)
	if err != nil {
		log.Fatalf("error loading environment variable: %v\n", err)
	}
	return value
}

func IntOrDefault(key string, defaultValue int) int {
	value, err := stringFromEnv(key)
	if err != nil {
		return defaultValue
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return number
}

func IntOrFatal(key string) int {
	value, err := stringFromEnv(key)
	if err != nil {
		log.Fatalf("error loading environment variable: %v\n", err)
	}
	number, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("error loading environment variable: %v\n", err)
	}
	return number
}

func stringFromEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", ErrMissingValue
}
