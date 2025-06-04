# Variables
DOCKER_COMPOSE = docker compose
IMAGE_NAME = kiononon2/lms-app
VERSION ?= latest

# Default target
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common targets:"
	@echo "  up            Start all services in detached mode"
	@echo "  down          Stop and remove all containers"
	@echo "  build         Build/rebuild all services"
	@echo "  restart       Restart all services"
	@echo "  logs          Show logs for all services"
	@echo "  ps            Show running containers"
	@echo "  rebuild       Stop, rebuild, and restart services"
	@echo "  exec          Open a shell in the app container"
	@echo ""
	@echo "Docker Hub targets:"
	@echo "  tag           Tag the built image with VERSION (default: latest)"
	@echo "  push          Push the image to Docker Hub with VERSION"
	@echo "  login         Login to Docker Hub"

.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: build
build:
	$(DOCKER_COMPOSE) build

.PHONY: restart
restart:
	$(DOCKER_COMPOSE) restart

.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

.PHONY: ps
ps:
	$(DOCKER_COMPOSE) ps

.PHONY: rebuild
rebuild: down build up

.PHONY: exec
exec:
	$(DOCKER_COMPOSE) exec app sh

.PHONY: tag
tag:
	docker tag $(IMAGE_NAME):latest $(IMAGE_NAME):$(VERSION)

.PHONY: push
push: tag
	docker push $(IMAGE_NAME):$(VERSION)

.PHONY: login
login:
	docker login
