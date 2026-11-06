SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

SWAGGER_PORT ?= 8081
FORMAT_SCOPE ?= all
DB_SERVICE ?= all
SANITY_SCOPE ?= ./tests/sanity/...
export SWAGGER_PORT

.PHONY: help
help:
	@bash ./scripts/help.sh

.PHONY: backend-setup backend-build backend-start backend-stop backend-status backend-dev backend-format
backend-setup:
	@bash ./scripts/backend/setup.sh

backend-build:
	@bash ./scripts/backend/build.sh

backend-start:
	@bash ./scripts/backend/start.sh

backend-stop:
	@bash ./scripts/backend/stop.sh

backend-status:
	@bash ./scripts/backend/status.sh

backend-dev:
	@bash ./scripts/backend/dev.sh

backend-format:
	@bash ./scripts/backend/format.sh "$(FORMAT_SCOPE)"

.PHONY: db-schema-status db-schema-init db-schema-reset
db-schema-status:
	@bash ./scripts/database/migrate.sh status "$(DB_SERVICE)"

db-schema-init:
	@bash ./scripts/database/migrate.sh init "$(DB_SERVICE)"

db-schema-reset:
	@bash ./scripts/database/migrate.sh reset "$(DB_SERVICE)"

.PHONY: nats-start nats-stop nats-status
nats-start:
	@bash ./scripts/infrastructure/nats/start.sh

nats-stop:
	@bash ./scripts/infrastructure/nats/stop.sh

nats-status:
	@bash ./scripts/infrastructure/nats/status.sh

.PHONY: sanity
sanity:
	@bash ./scripts/testing/sanity.sh "$(SANITY_SCOPE)"

.PHONY: swagger-start swagger-stop swagger-status
swagger-start:
	@bash ./scripts/swagger/start.sh

swagger-stop:
	@bash ./scripts/swagger/stop.sh

swagger-status:
	@bash ./scripts/swagger/status.sh

.PHONY: dev-start dev-start-with-logs dev-stop dev-status
dev-start:
	@bash ./scripts/dev/start.sh

dev-start-with-logs:
	@bash ./scripts/dev/start-with-logs.sh

dev-stop:
	@bash ./scripts/dev/stop.sh

dev-status:
	@bash ./scripts/dev/status.sh
