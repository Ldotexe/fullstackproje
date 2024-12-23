include .env

ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=${POSTGRES_USER} password=${POSTGRES_PW} dbname=${POSTGRES_DB} host=${POSTGRES_HOST} port=${POSTGRES_PORT} sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: run
run:
	go run ./cmd/service

.PHONY: run-all
run-all: up migration-up test run

