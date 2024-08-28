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

.PHONY: dc-up-pg
dc-up-pg:
	docker compose -f $(DOCKER_DIR)/docker-compose-pg.yml -p $(PROJECT_NAME) up --build

.PHONY: dc-up-mg
dc-up-mg:
	docker compose -f $(DOCKER_DIR)/docker-compose-mg.yml -p $(PROJECT_NAME) up --build	

.PHONY: dc-down-pg
dc-down-pg:
	docker compose -f $(DOCKER_DIR)/docker-compose-pg.yml -p $(PROJECT_NAME) down

.PHONY: dc-down-mg
dc-down-mg:
	docker compose -f $(DOCKER_DIR)/docker-compose-mg.yml -p $(PROJECT_NAME) down

.PHONY: docker-build-local-image
docker-local: ## build local image
	docker image build . -t $(PROJECT_NAME) -f ./docker/Dockerfile

.PHONY: lint
lint:
	golangci-lint run ./... --verbose --config=.golangci.yml

.PHONY: lint-fast
lint-fast:
	golangci-lint run ./... --fast --verbose --config=.golangci.yml

.PHONY: generate-structs
generate-structs:
	protoc -I ./internal/api/grpc/proto --go_out=./internal/api/grpc/gen/go/ --go_opt=paths=source_relative ./internal/api/grpc/proto/tpes/task_handler.proto --go-grpc_out=./internal/api/grpc/gen/go/ --go-grpc_opt=paths=source_relative