package config

import (
	"fmt"
)

// PostgresConfig is psql Config params
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	Password string `json:"password"`
}

// Dialect Database
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionInfo is about the information of db host, port, user, dbname
func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			c.Host,
			c.Port,
			c.User,
			c.DBName,
		)
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
	)
}

// DefaultPostgresConfig is Create psql config
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "lee",
		Password: "",
		DBName:   "support_dev",
	}
}

// GetSentryDSN is return the DNS info
func GetSentryDSN() string {
	return "https://61064a8e577b448ab6ed20f5aee63a1d@sentry.io/1777983"
}
