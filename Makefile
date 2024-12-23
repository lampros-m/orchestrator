BIN = ./bin
CMD_ORCHESTRATORSERVER = ./cmd/orchestratorserver
BIN_ORCHESTRATORSERVER = $(BIN)/orchestratorserver
CMD_ORCHESTRATORCLI = ./cmd/orchestratorcli
BIN_ORCHESTRATORCLI = $(BIN)/orchestratorcli

# Server build
.PHONY: build
build:
	CGO_ENABLED=0 go build -o "${BIN_ORCHESTRATORSERVER}" "${CMD_ORCHESTRATORSERVER}/"
	CGO_ENABLED=0 go build -o "${BIN_ORCHESTRATORCLI}" "${CMD_ORCHESTRATORCLI}/"

# Server orchestrator run
.PHONY: run-server
run-server:
	./"${BIN_ORCHESTRATORSERVER}"

# CLI orchestrator run
.PHONY: run-cli
run-cli:
	./"${BIN_ORCHESTRATORCLI}"

# Swagger
# go install github.com/swaggo/swag/cmd/swag@latest
.PHONY: swag
swag:
	swag init -g internal/apihttp/router.go
	swag fmt