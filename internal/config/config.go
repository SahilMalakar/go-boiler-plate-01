package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseURL string
	HTTP        HTTPConfig // HTTP server configuration.
	DB          DBConfig   // Database configuration.
}

type HTTPConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	PingTimeout     time.Duration
}

func MustLoad() *Config {
	// Loads variables from the .env file into the environment.
	godotenv.Load()

	return &Config{
		Port:        getStringEnv("PORT"),
		Env:         getStringEnv("ENV"),
		DatabaseURL: getStringEnv("DATABASE_URL"),

		HTTP: HTTPConfig{
			ReadTimeout:  getDurationEnv("HTTP_READ_TIMEOUT"),
			WriteTimeout: getDurationEnv("HTTP_WRITE_TIMEOUT"),
			IdleTimeout:  getDurationEnv("HTTP_IDLE_TIMEOUT"),
		},

		DB: DBConfig{
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: time.Duration(getIntEnv("DB_CONN_MAX_LIFETIME")) * time.Minute,
			PingTimeout:     getDurationEnv("DB_PING_TIMEOUT"),
		},
	}
}

// Reads a required string environment variable.
func getStringEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		panic(key + " is required")
	}

	return value
}

// Reads a required integer environment variable.
func getIntEnv(key string) int {
	value := os.Getenv(key)

	if value == "" {
		panic(key + " is required")
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		panic(key + " must be a valid integer")
	}

	return number
}

// Reads a required timeout value in seconds
// and converts it into time.Duration.
func getDurationEnv(key string) time.Duration {
	value := getIntEnv(key)

	return time.Duration(value) * time.Second
}
