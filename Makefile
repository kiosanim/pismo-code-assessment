.PHONY: go docker migration_tool

MIGRATION_TOOL = "cmd/migrate/migration_tool.go"

#create.config:
#	@echo "TODO: CRIAR CONFIG FILE" <-------------------- FALTA DESENVOLVER

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
	go test -v ./...

docker-up:
	@echo "Provisioning all containers from docker-compose.yaml"
	docker compose up -d

docker-destroy:
	@echo "Destroying all containers from docker-compose.yaml removing orphans and volumes"
	docker compose down --remove-orphans --volumes

docker-up-postgres:
	@echo "Initializing database container"
	docker compose up -d postgres

docker-down-db:
	@echo "Stop and remove database container"
	docker compose down -d postgres

