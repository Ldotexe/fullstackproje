#!/usr/bin/env bash

SETUP_TEST="user=${POSTGRES_USER} password=${POSTGRES_PW} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable"
INTERNAL_PKG_PATH="/base"
MIGRATION_FOLDER="${INTERNAL_PKG_PATH}/migrations"

#goose -dir "${MIGRATION_FOLDER}" create "${POSTGRES_DB}" sql
goose -dir "${MIGRATION_FOLDER}" postgres "user=${POSTGRES_USER} password=${POSTGRES_PW} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable" up
#goose -dir "${MIGRATION_FOLDER}" postgres "${POSTGRES_SETUP_TEST}" down

go run ./cmd/service
