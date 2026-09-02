package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Environment   string
	AllowedHosts  string
	LabPassword   string
	SessionSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("ENV", "development"),
		AllowedHosts:  getEnv("ALLOWED_HOSTS", ""),
		LabPassword:   getEnv("LAB_PASSWORD", "changeme"),
		SessionSecret: getEnv("SESSION_SECRET", "dev-secret-32-chars-minimum-here"),
	}

	return cfg, cfg.Validate()
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) Validate() error {
	var errs []string

	switch c.Environment {
	case "development", "production":
	default:
		errs = append(errs, fmt.Sprintf("ENV must be 'development' or 'production', got %q", c.Environment))
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Sprintf("PORT must be a valid port number, got %q", c.Port))
	}

	if c.IsProduction() && c.AllowedHosts == "" {
		errs = append(errs, "ALLOWED_HOSTS is required in production")
	}

	if c.IsProduction() && c.LabPassword == "changeme" {
		errs = append(errs, "LAB_PASSWORD must be set in production")
	}

	if len(c.SessionSecret) < 32 {
		errs = append(errs, "SESSION_SECRET must be at least 32 characters")
	}

	if c.AllowedHosts != "" {
		for _, host := range strings.Split(c.AllowedHosts, ",") {
			host = strings.TrimSpace(host)
			if host == "" {
				continue
			}
			if strings.Contains(host, "://") {
				errs = append(errs, fmt.Sprintf("ALLOWED_HOSTS must not include scheme, got %q", host))
			}
			h, _, err := net.SplitHostPort(host)
			if err != nil {
				h = host
			}
			if h == "" {
				errs = append(errs, fmt.Sprintf("ALLOWED_HOSTS contains invalid entry: %q", host))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("configuration errors:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
