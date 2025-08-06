ENV_FILE := .env

include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

.PHONY: up down restart rebuild clean-db seed-db

DOCKER_COMPOSE_FILE := docker-compose.yaml
SEEDER_DIR := ./internal/adapters/database/seeder
POSTGRES_DB := $(POSTGRES_DB)
POSTGRES_USER := $(POSTGRES_USER)
POSTGRES_PASSWORD := $(POSTGRES_PASSWORD)

# =================================================================================================
# Docker Commands
# =================================================================================================

up: ## Start the containers in detached mode.
	@echo "Starting containers..."
	docker-compose --env-file $(ENV_FILE) up -d --force-recreate --remove-orphans
	sleep 3
	go run ./cmd/app/main.go ./cmd/app/bootstrap.go &

down: ## Stop and remove the containers, networks, and volumes.
	@echo "Stopping and removing containers..."
	docker-compose --env-file $(ENV_FILE) down --volumes --remove-orphans
	$(MAKE) kill-go || true

restart: ## Restart all containers.
	$(MAKE) down
	$(MAKE) up

# =================================================================================================
# Database Commands
# =================================================================================================

seed-db: ## Run SQL seeder files against the database.
	@echo "Running database seeders..."
	@for file in $(SEEDER_DIR)/*.sql; do \
		if [ -f "$$file" ]; then \
			echo "Executing $$file..."; \
			docker exec -i postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < "$$file"; \
		fi \
	done
	@echo "Database seeding complete."

# =================================================================================================
# App Commands
# =================================================================================================

up-and-seed: ## Restart containers, clean database, and run seeders.
	$(MAKE) restart
	sleep 2
	$(MAKE) seed-db

kill-go: ## Kill the running Go process.
	@echo "Killing Go process on port $(APP_PORT)..."
	-lsof -t -i:$(APP_PORT) | xargs kill || true


# =================================================================================================
# Helper Commands
# =================================================================================================

help: ## Show this help message.
	@echo "Usage: make [command]"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?##/ {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)