package config

import (
	"os"
	"strconv"
	"sync"
)

const (
	DefaultHTTPServerPort      = "8080"
	DefaultLogLevel            = "info"
	DefaultPostgresPort        = "5432"
	DefaultPostgresSSLMode     = "disable"
	DefaultPostgresMaxPoolSize = 5
	DefaultPostgresMinPoolSize = 1
)

var (
	instance *Config
	once     sync.Once
)

// Config is a configuration struct of service
type Config struct {
	HTTPServerPort      string
	LogLevel            string
	PostgresHost        string
	PostgresPort        string
	PostgresDatabase    string
	PostgresUsername    string
	PostgresPassword    string
	PostgresSSLMode     string
	PostgresMaxPoolSize int
	PostgresMinPoolSize int
	PostgresLogEnabled  bool
}

// GetConfig returns the configuration from environment variables
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{

			HTTPServerPort:      getDefaultEnv("HTTP_SERVER_PORT", DefaultHTTPServerPort),
			LogLevel:            getDefaultEnv("LOG_LEVEL", DefaultLogLevel),
			PostgresHost:        os.Getenv("POSTGRES_HOST"),
			PostgresPort:        getDefaultEnv("POSTGRES_PORT", DefaultPostgresPort),
			PostgresDatabase:    os.Getenv("POSTGRES_DATABASE"),
			PostgresUsername:    os.Getenv("POSTGRES_USERNAME"),
			PostgresPassword:    os.Getenv("POSTGRES_PASSWORD"),
			PostgresSSLMode:     getDefaultEnv("POSTGRES_SSLMODE", DefaultPostgresSSLMode),
			PostgresMaxPoolSize: getDefaultIntEnv("POSTGRES_MAX_POOL_SIZE", DefaultPostgresMaxPoolSize),
			PostgresMinPoolSize: getDefaultIntEnv("POSTGRES_MIN_POOL_SIZE", DefaultPostgresMinPoolSize),
			PostgresLogEnabled:  getBooleanFromValue(os.Getenv("POSTGRES_LOG_ENABLED")),
		}
	})
	return instance
}

func getDefaultEnv(name string, value string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	return value
}

func getDefaultIntEnv(name string, defaultValue int) int {
	if s := os.Getenv(name); s != "" {
		v, err := strconv.Atoi(s)
		if err == nil {
			return v
		}
	}
	return defaultValue
}

func getBooleanFromValue(value string) bool {
	boolean, err := strconv.ParseBool(value)

	if err != nil {
		return false
	}

	return boolean
}
