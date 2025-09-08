# Read .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Variables
dev:
	docker compose -f $(DOCKER_COMPOSE_FILE_DEV) up --build
