package envconfig

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadIfDev loads .env file if APP_ENV != "production"
func LoadIfDev() {
	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load()
	}
}

// IsProd returns true if APP_ENV == "production"
func IsProd() bool {
	return os.Getenv("APP_ENV") == "production"
}

// String returns the value of an environment variable,
// or the provided defaultValue if it's empty.
func String(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// Int returns the value of an environment variable as int,
// or the provided defaultValue if it's empty or invalid.
func Int(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return defaultValue
}
