SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

SWAGGER_PORT ?= 8081
export SWAGGER_PORT

.PHONY: help
help:
	@./scripts/help.sh

.PHONY: backend-build backend-start backend-stop backend-status backend-dev
backend-build:
	@./scripts/backend/build.sh

backend-start:
	@./scripts/backend/start.sh

backend-stop:
	@./scripts/backend/stop.sh

backend-status:
	@./scripts/backend/status.sh

backend-dev:
	@./scripts/backend/dev.sh

.PHONY: swagger-start swagger-stop swagger-status
swagger-start:
	@./scripts/swagger/start.sh

swagger-stop:
	@./scripts/swagger/stop.sh

swagger-status:
	@./scripts/swagger/status.sh

.PHONY: dev-start dev-stop dev-status
dev-start:
	@./scripts/dev/start.sh

dev-stop:
	@./scripts/dev/stop.sh

dev-status:
	@./scripts/dev/status.sh
