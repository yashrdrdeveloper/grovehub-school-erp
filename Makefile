SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

SWAGGER_PORT ?= 8081
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

.PHONY: swagger-start swagger-stop swagger-status
swagger-start:
	@bash ./scripts/swagger/start.sh

swagger-stop:
	@bash ./scripts/swagger/stop.sh

swagger-status:
	@bash ./scripts/swagger/status.sh

.PHONY: dev-start dev-stop dev-status
dev-start:
	@bash ./scripts/dev/start.sh

dev-stop:
	@bash ./scripts/dev/stop.sh

dev-status:
	@bash ./scripts/dev/status.sh
