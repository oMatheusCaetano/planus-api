# Read .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Variables
DOCKER_EXEC_API = docker exec -it $(API_CONTAINER_NAME)
MIGRATIONS_DIR = database/migrations
DB_CONNECTION_STRING = "user=$(DB_USER) password=$(DB_PASSWORD) host=$(DB_HOST) port=$(DB_CONTAINER_PORT) dbname=$(DB_NAME) sslmode=disable"

# Commands
dev:
	docker compose -f $(DOCKER_COMPOSE_FILE_DEV) up --build

dev-detached:
	docker compose -f $(DOCKER_COMPOSE_FILE_DEV) up -d --build

dev-down:
	docker compose -f $(DOCKER_COMPOSE_FILE_DEV) down

get:
	$(DOCKER_EXEC_API) go get $(filter-out $@,$(MAKECMDGOALS))

migration-create:
	$(DOCKER_EXEC_API) goose create -dir $(MIGRATIONS_DIR) $(filter-out $@,$(MAKECMDGOALS)) sql

migration-up:
	$(DOCKER_EXEC_API) goose -dir $(MIGRATIONS_DIR) postgres $(DB_CONNECTION_STRING) up

seed:
	$(DOCKER_EXEC_API) go run cmd/seed/main.go

%:
	@:
