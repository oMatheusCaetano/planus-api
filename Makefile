ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Variables
DOCKER_EXEC_API = docker exec -it $(API_CONTAINER_NAME)

# Commands
run:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build --force-recreate

test:
	@$(DOCKER_EXEC_API) go test ./... -covermode=atomic

tidy:
	@$(DOCKER_EXEC_API) go mod tidy
