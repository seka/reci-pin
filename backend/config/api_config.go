package config

import (
	"fmt"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Storage  StorageConfig
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

type StorageConfig struct {
	Bucket        string
	Endpoint      string // Optional, for LocalStack
	PublicBaseURL string // Optional, for public access URL
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
