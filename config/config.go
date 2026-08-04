package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv loads .env.dv or .env.pd based on APP_ENV (default: dv).
// APP_ENV can be set in the process environment before startup.
func LoadEnv() error {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" {
		env = "dv"
	}

	filename := ".env." + env
	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("failed to load %s: %w", filename, err)
	}
	return nil
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func AuthToken() string {
	return os.Getenv("AUTH_TOKEN")
}

func UseMockDB() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("USE_MOCK_DB")))
	return v == "true" || v == "1" || v == "yes"
}

func DBURL() string {
	return os.Getenv("DB_URL")
}

func CORSOrigins() []string {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		return []string{"http://localhost:3000"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if o := strings.TrimSpace(p); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
