.PHONY: build test start migrate bootstrap-operator
build:
	GOWORK=off go build ./...
test:
	GOWORK=off go test ./...
start:
	GOWORK=off go run ./cmd/internal-tool

migrate:
	GOWORK=off go run -mod=mod ./cmd/operator-admin migrate
bootstrap-operator:
	@GOWORK=off go run -mod=mod ./cmd/operator-admin bootstrap "$(EMAIL)"
