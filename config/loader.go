package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var MainConfig *Config

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, reading from environment")
	}

	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	validate(cfg)

	MainConfig = cfg
	return cfg
}

func validate(cfg *Config) {
	var missing []string

	// App
	if cfg.App.Name == "" {
		missing = append(missing, "APP_NAME")
	}
	if cfg.App.Env == "" {
		missing = append(missing, "APP_ENV")
	}
	if cfg.App.Port == "" {
		missing = append(missing, "APP_PORT")
	}
	if cfg.App.RequestTimeout == 0 {
		missing = append(missing, "APP_REQUESTTIMEOUT")
	}
	if cfg.App.ReadTimeout == 0 {
		missing = append(missing, "APP_READTIMEOUT")
	}
	if cfg.App.WriteTimeout == 0 {
		missing = append(missing, "APP_WRITETIMEOUT")
	}
	if cfg.App.IdleTimeout == 0 {
		missing = append(missing, "APP_IDLETIMEOUT")
	}

	// Database
	if cfg.Database.Host == "" {
		missing = append(missing, "DATABASE_HOST")
	}
	if cfg.Database.Port == "" {
		missing = append(missing, "DATABASE_PORT")
	}
	if cfg.Database.User == "" {
		missing = append(missing, "DATABASE_USER")
	}
	if cfg.Database.Password == "" {
		missing = append(missing, "DATABASE_PASSWORD")
	}
	if cfg.Database.Name == "" {
		missing = append(missing, "DATABASE_NAME")
	}
	if cfg.Database.MaxIdleConn == 0 {
		missing = append(missing, "DATABASE_MAXIDLECONN")
	}
	if cfg.Database.SSLMode == "" {
		missing = append(missing, "DATABASE_SSLMODE")
	}

	// JWT
	if cfg.JWT.SecretKey == "" {
		missing = append(missing, "JWT_SECRETKEY")
	}
	if cfg.JWT.ExpireDuration == 0 {
		missing = append(missing, "JWT_EXPIREDURATION")
	}

	if len(missing) > 0 {
		log.Fatalf(
			"\n[CONFIG] Missing or zero-value environment variables (%d):\n  - %s\n",
			len(missing),
			strings.Join(missing, "\n  - "),
		)
	}
}
