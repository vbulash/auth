#!/bin/bash
source .env

export MIGRATION_DSN="host=postgres port=5432 dbname=$DB_DATABASE user=$DB_USER password=$DB_PASSWORD sslmode=disable"
goose -dir . postgres "$MIGRATION_DSN" up -v
