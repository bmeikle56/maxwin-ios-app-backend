package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// LoadEnv loads .env.dv or .env.pd based on APP_ENV (default: dv).
// On Railway, vars come from the platform — missing files are OK if JWT_SECRET is set.
func LoadEnv() error {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if env == "" {
		env = "dv"
		_ = os.Setenv("APP_ENV", env)
	}

	filename := ".env." + env
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			if strings.TrimSpace(os.Getenv("JWT_SECRET")) == "" {
				return fmt.Errorf("%s not found and JWT_SECRET is not set", filename)
			}
			return nil
		}
		return fmt.Errorf("stat %s: %w", filename, err)
	}

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
