package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig_DSN(t *testing.T) {
	c := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "password",
		DBName:   "mydb",
		SSLMode:  "disable",
	}

	dsn := c.DSN()
	t.Logf("Generated DSN: %s", dsn)
	// Expected format: postgres://user:password@localhost:5432/mydb?sslmode=disable
	assert.Contains(t, dsn, "postgres://user:password@localhost:5432/mydb")
	assert.Contains(t, dsn, "sslmode=disable")

	// Case: Empty DBName
	c.DBName = ""
	dsnEmpty := c.DSN()
	t.Logf("Generated DSN (empty DB): %s", dsnEmpty)
	assert.Contains(t, dsnEmpty, "postgres://user:password@localhost:5432")
	assert.NotContains(t, dsnEmpty, "/?")
}
