package config

import (
	"net"
	"net/url"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

type StorageConfig struct {
	Bucket        string
	Endpoint      string // Optional, for LocalStack
	PublicBaseURL string // Optional, for public access URL
}

type JWTConfig struct {
	Secret                     string
	ExpirationHours            int
	RefreshTokenExpirationDays int
}

func (c *DatabaseConfig) DSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
