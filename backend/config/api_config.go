package config

import (
	"net"
	"net/url"
)

type Config struct {
	Database  Database
	ApiServer ApiServer
	Storage   Storage
}

type ApiServer struct {
	Port string
	JWT  JWT
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Storage struct {
	Bucket        string
	Endpoint      string // Optional, for LocalStack
	PublicBaseURL string // Optional, for public access URL
}

type JWT struct {
	Secret                     string
	ExpirationHours            int
	RefreshTokenExpirationDays int
}

func (c *Database) DSN() string {
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
