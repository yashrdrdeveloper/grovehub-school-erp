SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

SWAGGER_PORT ?= 8081
FORMAT_SCOPE ?= all
LINT_SCOPE ?= foundation-repository
DB_SERVICE ?= all
SANITY_SCOPE ?= ./tests/sanity/...
FRONTEND_TUNNEL ?= 1
export SWAGGER_PORT FRONTEND_TUNNEL

.PHONY: help
help:
	@bash ./scripts/help.sh

.PHONY: backend-setup backend-build backend-test backend-start backend-stop backend-status backend-dev backend-format backend-lint
backend-setup:
	@bash ./scripts/backend/setup.sh

backend-build:
	@bash ./scripts/backend/build.sh

backend-test:
	@bash ./scripts/backend/test.sh ./...

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

backend-lint:
	@bash ./scripts/backend/format.sh "$(LINT_SCOPE)" --check
	@bash ./scripts/backend/check-identity-architecture.sh

.PHONY: db-schema-status db-schema-init db-schema-reset
db-schema-status:
	@bash ./scripts/database/migrate.sh status "$(DB_SERVICE)"

db-schema-init:
	@bash ./scripts/database/migrate.sh init "$(DB_SERVICE)"

db-schema-reset:
	@bash ./scripts/database/migrate.sh reset "$(DB_SERVICE)"

# package=foundation or packages='["foundation","identity"]'; omitted = all.
.PHONY: unit-test unit-coverage test-inventory
unit-test:
	@UNIT_PACKAGE='$(package)' UNIT_PACKAGES='$(packages)' bash ./scripts/testing/unit-coverage.sh --report-only

unit-coverage:
	@UNIT_PACKAGE='$(package)' UNIT_PACKAGES='$(packages)' bash ./scripts/testing/unit-coverage.sh --strict

test-inventory:
	@UNIT_PACKAGE='$(package)' UNIT_PACKAGES='$(packages)' bash ./scripts/testing/unit-coverage.sh --inventory

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

.PHONY: frontend-setup frontend-start frontend-stop frontend-status frontend-build frontend-test frontend-dev
frontend-setup:
	@bash ./scripts/frontend/setup.sh

frontend-start:
	@bash ./scripts/frontend/start.sh

frontend-stop:
	@bash ./scripts/frontend/stop.sh

frontend-status:
	@bash ./scripts/frontend/status.sh

frontend-build:
	@bash ./scripts/frontend/build.sh

frontend-test:
	@cd frontend/projects/school && pnpm test

frontend-dev:
	@cd frontend/projects/school && pnpm dev

# Operator-only, one-time first account bootstrap. Prompts for password securely.
.PHONY: identity-bootstrap-admin
identity-bootstrap-admin:
	@bash ./scripts/identity/bootstrap-admin.sh "$(EMAIL)" "$(FIRST_NAME)" "$(LAST_NAME)" "$(TENANT_ID)" $(if $(PERSON_GID),"$(PERSON_GID)")

# Explicit project entrypoints. Each command uses its own environment.
.PHONY: school-backend-build school-backend-test school-backend-start internal-backend-build internal-backend-test internal-backend-start internal-db-migrate internal-frontend-build internal-frontend-test internal-frontend-dev
school-backend-build:
	@$(MAKE) -C backend/projects/school build
school-backend-test:
	@$(MAKE) -C backend/projects/school test
school-backend-start:
	@$(MAKE) -C backend/projects/school start
internal-backend-build:
	@$(MAKE) -C backend/projects/internal-tool build
internal-backend-test:
	@$(MAKE) -C backend/projects/internal-tool test
internal-backend-start:
	@$(MAKE) -C backend/projects/internal-tool start
internal-db-migrate:
	@$(MAKE) -C backend/projects/internal-tool migrate
internal-frontend-build:
	@cd frontend && pnpm --filter @groovehub/internal-tool build
internal-frontend-test:
	@cd frontend && pnpm --filter @groovehub/internal-tool test
internal-frontend-dev:
	@cd frontend && pnpm --filter @groovehub/internal-tool dev

.PHONY: platform-lib-test platform-lib-coverage
platform-lib-test:
	@cd backend/libs/platform && GOWORK=off go test ./...

platform-lib-coverage:
	@bash ./backend/libs/platform/tools/coverage/check.sh
