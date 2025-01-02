BIN = ./bin
CMD_ORCHESTRATORSERVER = ./cmd/orchestratorserver
BIN_ORCHESTRATORSERVER = $(BIN)/orchestratorserver

# Server build
.PHONY: build
build:
	CGO_ENABLED=0 go build -o "${BIN_ORCHESTRATORSERVER}" "${CMD_ORCHESTRATORSERVER}"

# Server orchestrator run
.PHONY: run-server
run-server:
	./"${BIN_ORCHESTRATORSERVER}"

# Swagger
# go install github.com/swaggo/swag/cmd/swag@latest
.PHONY: swag
swag:
	swag init -g internal/apihttp/router.go
	swag fmt