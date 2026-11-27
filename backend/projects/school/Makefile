.PHONY: generate-clients build test start
generate-clients:
	bash scripts/generate-clients.sh
build: generate-clients
	GOWORK=off go build ./...
test: generate-clients
	GOWORK=off go test ./...
start: generate-clients
	GOWORK=off go run ./cmd/school
