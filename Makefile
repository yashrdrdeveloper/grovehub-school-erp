SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

SWAGGER_PORT ?= 8081
DB_SERVICE ?= all
export SWAGGER_PORT

.PHONY: help
help:
	@bash ./scripts/help.sh

.PHONY: backend-setup backend-build backend-start backend-stop backend-status backend-dev
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

.PHONY: mail-provider-start mail-provider-stop mail-provider-status mail-domain-ready
mail-provider-start:
	@bash ./scripts/mailProvider/start.sh

mail-provider-stop:
	@bash ./scripts/mailProvider/stop.sh

mail-provider-status:
	@bash ./scripts/mailProvider/status.sh

mail-domain-ready:
	@bash ./scripts/mailProvider/domain-ready.sh

.PHONY: db-schema-status db-schema-init db-schema-reset
db-schema-status:
	@bash ./scripts/database/migrate.sh status "$(DB_SERVICE)"

db-schema-init:
	@bash ./scripts/database/migrate.sh init "$(DB_SERVICE)"

db-schema-reset:
	@bash ./scripts/database/migrate.sh reset "$(DB_SERVICE)"

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
