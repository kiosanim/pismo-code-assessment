.PHONY: go docker migration_tool swag

MIGRATION_TOOL = "cmd/migrate/migration_tool.go"
IMAGE_NAME = pismo-code-assesment

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
	go run $(MIGRATION_TOOL) up

db-migration-down:
	@echo "Running Migration Tool down"
	go run $(MIGRATION_TOOL) down

db-migration-status:
	@echo "Running Migration Tool status"
	go run $(MIGRATION_TOOL) status

test-unit:
	@echo "Running unit tests"
	go test -v ./...

docker-build:
	docker build

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
