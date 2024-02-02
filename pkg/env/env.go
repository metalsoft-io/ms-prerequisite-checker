package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

func GetString(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("Invalid value '%s' for variable '%s'. Using the default value %d.\n", value, key, fallback)
			return fallback
		}
		return valueInt
	}
	return fallback
}

func GetFloat64(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		valueFloat64, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Invalid value '%s' for variable '%s'. Using the default value %v.\n", value, key, fallback)
			return fallback
		}
		return valueFloat64
	}
	return fallback
}

func GetBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		valueBool, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Printf("Invalid value '%s' for variable '%s'. Using the default value %v.\n", value, key, fallback)
			return fallback
		}
		return valueBool
	}
	return fallback
}

func ParseLogLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		fmt.Printf("Invalid log level value '%s'. Using the default log level INFO.\n", logLevel)
		return zerolog.InfoLevel
	}
}
