package config

import (
	"fmt"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	DBName   string `json:"dbname"`
	Password string `json:"password"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

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

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "lee",
		Password: "",
		DBName:   "support_dev",
	}
}
