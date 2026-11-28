.PHONY: generate-clients build test start migrate bootstrap-operator
generate-clients:
	bash scripts/generate-clients.sh
build: generate-clients
	GOWORK=off go build -mod=mod ./...
test: generate-clients
	GOWORK=off go test -mod=mod ./...
start: generate-clients
	GOWORK=off go run -mod=mod ./cmd/internal-tool

migrate:
	GOWORK=off go run -mod=mod ./cmd/operator-admin migrate
bootstrap-operator:
	@GOWORK=off go run -mod=mod ./cmd/operator-admin bootstrap "$(EMAIL)"
