# Read .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Variables
DOCKER_EXEC_API = docker exec -it docker-compose.dev.yml

# Commands
dev:
	docker compose -f $(DOCKER_COMPOSE_FILE_DEV) up --build

get:
	$(DOCKER_EXEC_API) go get $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
