package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	App    AppConfig
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type ServerConfig struct {
	port            int
	ShutdownTimeout time.Duration
}

type AppConfig struct {
	Env string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}
	return &Config{
		DB: DBConfig{
			DSN:          getEnv("FOLLOW_DB_DSN", "postgres://greenlight:your_password@localhost:5432/follow?sslmode=disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  getEnvAsDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		Server: ServerConfig{
			port:            port,
			ShutdownTimeout: 5 * time.Second,
		},
		App: AppConfig{
			Env: os.Getenv("APP_ENV"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		i, _ := strconv.Atoi(val)
		return i
	}
	return defaultValue
}

func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	if valueStr, exists := os.LookupEnv(name); exists {
		value, err := time.ParseDuration(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultVal
}
