ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Variables
DOCKER_EXEC_API = docker exec -it $(API_CONTAINER_NAME)

# Commands
run:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build --force-recreate

run.detach:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build --force-recreate --detach

test:
	@$(DOCKER_EXEC_API) go test ./... -covermode=atomic

tidy:
	@$(DOCKER_EXEC_API) go mod tidy

get:
	@$(DOCKER_EXEC_API) go get $(filter-out $@,$(MAKECMDGOALS))

swagger:
	@$(DOCKER_EXEC_API) swag init -g cmd/api/main.go -o docs
