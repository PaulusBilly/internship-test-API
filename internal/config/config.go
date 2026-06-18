package config

import (
	"os"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func Load() *Config {
	return &Config{
		Host:     envOrDefault("DB_HOST", "localhost"),
		Port:     envOrDefault("DB_PORT", "5432"),
		Username: envOrDefault("DB_USER", "postgres"),
		Password: envOrDefault("DB_PASS", "postgres"),
		DBName:   envOrDefault("DB_NAME", "viewdata"),
	}
}

func (c *Config) DBConnStr() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.Password, c.DBName,
	)
}

func envOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
