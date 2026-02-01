package config

import (
	"fmt"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port int
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
