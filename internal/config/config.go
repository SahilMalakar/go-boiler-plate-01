package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func MustLoad() *Config {
	// Loads variables from the .env file into the environment.
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	// Reads and converts Timeout from string to integer.
	readTimeout := getIntEnv("ReadTimeout")
	writeTimeout := getIntEnv("WriteTimeout")
	idleTimeout := getIntEnv("IdleTimeout")

	return &Config{
		Port:         port,
		Env:          env,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}
}

// Reads an environment variable and converts it from string to int.
func getIntEnv(key string) int {
	value := os.Getenv(key)

	if value == "" {
		panic(key + " is required")
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		panic(key + " must be a valid number")
	}

	return number
}
