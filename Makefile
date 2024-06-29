SRC_DIR := ./cmd/tpes
MAIN_FILE := $(SRC_DIR)/main.go

BINARY_NAME := ./bin/tpes

DOCKER_DIR := ./docker
PROJECT_NAME := tpes

.PHONY: build
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: docker-compose-up
docker-compose-up:
	docker compose -f $(DOCKER_DIR)/docker-compose.yml -p $(PROJECT_NAME) up --build

.PHONY: docker-compose-down
docker-compose-down:
	docker compose -f $(DOCKER_DIR)/docker-compose.yml -p $(PROJECT_NAME) down

.PHONY: docker-build-local-image
docker-local: ## build local image
	docker image build . -t $(PROJECT_NAME) -f ./docker/Dockerfile

## в чем смысл таких локальных образов бд, если каждый раз удаляются таблицы в базе?
.PHONY: docker-postgres-start 
docker-postgres-start: ## start the postgres container
	docker run --rm --name postgres	-e POSTGRES_PASSWORD=12345 -e POSTGRES_DB=postgres -d -p 5432:5432 postgres

.PHONY: docker-postgres-stop
docker-postgres-stop: ## stop the postgres container
	docker stop postgres

.PHONY: docker-tpes-start
docker-tpes-start: ## start the app from container
	docker run --name my-container tpes-app

.PHONY: lint
lint:
	golangci-lint run ./... --verbose --config=.golangci.yml

.PHONY: lint-fast
lint-fast:
	golangci-lint run ./... --fast --verbose --config=.golangci.yml