# Variables
default: help
APP_NAME := api
STUB_NAME := stub-server
DOCKER_COMPOSE := docker compose
GOBIN ?= $(CURDIR)/bin
.PHONY: help build-stub build-api build all up down logs

help:
	@echo "Usage: make [target]"
	@echo "  build-stub     Build Docker image for stub server"
	@echo "  build-api      Build Docker image for API service"
	@echo "  all            Build both images"
	@echo "  up             Start services via docker-compose"
	@echo "  down           Stop services"
	@echo "  logs           Tail logs of both services"

build-stub:
	docker build -f Dockerfile.stub -t $(STUB_NAME) .

build-api:
	docker build -f Dockerfile.api -t $(APP_NAME) .

build: build-stub build-api

up:
	$(DOCKER_COMPOSE) up --build -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

install-cli:
	@mkdir -p $(GOBIN)
	@go build -o $(GOBIN)/taskcli ./cmd/taskcli
	@echo "✅ taskCLI built in $(GOBIN)"

CLI: install-cli
	@echo "🚀 Running taskCLI shell"
	@PATH="$(GOBIN):$(PATH)" exec $(SHELL) -i

run-stub:
	go run $(CURDIR)/cmd/io_server/main.go
run-api:
	go run $(CURDIR)/cmd/api/main.go