.PHONY: go docker migration_tool swag

MIGRATION_TOOL = "cmd/migrate/migration_tool.go"
MIGRATION_TOOL_BIN = "/app/app_migration_tool"
IMAGE_NAME = pismo-api
CONTAINER_NAME = pismo-api

#create.config:
#	@echo "TODO: CRIAR CONFIG FILE" <-------------------- FALTA DESENVOLVER

run:
	@echo "Running API Server"
	go run cmd/api/main.go

install:
	@echo "Running Migration Tool create"
	go mod tidy
	go mod vendor

db-migration-up:
	@echo "Running Migration Tool create"
	go run cmd/migrate/main.go up

db-migration-down:
	@echo "Running Migration Tool"
	go run $(MIGRATION_TOOL) down

db-migration-status:
	@echo "Running Migration Tool status"
	go run $(MIGRATION_TOOL) status

db-migration-docker-up:
	@echo "Running Migration Tool create from docker container"
	docker exec -it $(CONTAINER_NAME) $(MIGRATION_TOOL_BIN) up

db-migration-docker-down:
	@echo "Running Migration Tool down from docker container"
	docker exec -it $(CONTAINER_NAME) $(MIGRATION_TOOL_BIN) down

db-migration-docker-status:
	@echo "Running Migration Tool status from docker container"
	docker exec -it $(CONTAINER_NAME) $(MIGRATION_TOOL_BIN) status

test-unit:
	@echo "Running unit tests"
	go test -v ./...

docker-provision: docker-up db-migration-docker-up

docker-up:
	@echo "Provisioning all containers from docker-compose.yaml"
	docker compose up -d --build

docker-destroy:
	@echo "Destroying all containers from docker-compose.yaml removing orphans and volumes"
	docker compose down --remove-orphans --volumes

docker-up-postgres:
	@echo "Initializing database container"
	docker compose up -d postgres

docker-down-db:
	@echo "Stop and remove database container"
	docker compose down -d postgres

swagger:
	@echo "Generating Swagger Docs"
	swag init -g cmd/api/main.go -o docs
