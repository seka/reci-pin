#!/bin/sh
set -e

# Construct command flags from environment variables
PORT=${BACKEND_PORT:-8080}
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-recipin_dev}
DB_SSLMODE=${DB_SSLMODE:-disable}
JWT_SECRET=${JWT_SECRET:-change-me}
JWT_EXPIRATION=${JWT_EXPIRATION_HOURS:-24}

# If the first argument is "air", we assume we are in dev mode.
# air invokes the binary itself, so we can't pass flags to air directly to affect the binary easily
# without modifying .air.toml.
# HOWEVER, for the final "production" run (which this entrypoint is also for), we execute the binary.

if [ "$1" = "air" ]; then
    # For development with air, we need to pass these arguments to the binary.
    # Air allows passing arguments to the binary after --
    exec air "$@" -- \
        -port="$PORT" \
        -db-host="$DB_HOST" \
        -db-port="$DB_PORT" \
        -db-user="$DB_USER" \
        -db-password="$DB_PASSWORD" \
        -db-name="$DB_NAME" \
        -db-sslmode="$DB_SSLMODE" \
        -jwt-secret="$JWT_SECRET" \
        -jwt-expiration="$JWT_EXPIRATION"
else
    # For normal execution (e.g. ./main)
    exec "$@" \
        -port="$PORT" \
        -db-host="$DB_HOST" \
        -db-port="$DB_PORT" \
        -db-user="$DB_USER" \
        -db-password="$DB_PASSWORD" \
        -db-name="$DB_NAME" \
        -db-sslmode="$DB_SSLMODE" \
        -jwt-secret="$JWT_SECRET" \
        -jwt-expiration="$JWT_EXPIRATION"
fi
