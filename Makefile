PROJECT_NAME=grpc-pb
MODULE_NAME=grpc-pb

.DEFAULT_GOAL := build

.PHONY: build
build:
	@go build ./cmd/grpc-pb/

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -coverprofile coverage.out ./...

.PHONY: coverage
coverage:
	@go tool cover -html=coverage.out

.PHONY: get
get:
	@go mod download



