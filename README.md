# Pismo Code Assessment - Complete Technical Documentation

**Author:** Fábio Sartori

---

## Quick Start

To run the complete application with Docker (including database and migrations):

```bash
make docker-provision
```

This single command will:
1. Build and start all Docker containers (API + PostgreSQL)
2. Run all database migrations automatically
3. Make the API available at `http://localhost:8080`
4. Swagger documentation at `http://localhost:8080/swagger/index.html`

**That's it!** The application is now running and ready to use.

---

## Table of Contents
1. [Quick Start](#quick-start)
2. [Project Overview](#project-overview)
3. [Architecture](#architecture)
4. [Technology Stack](#technology-stack)
5. [Project Structure](#project-structure)
6. [Domain Layer](#domain-layer)
7. [Application Layer](#application-layer)
8. [Infrastructure Layer](#infrastructure-layer)
9. [Interface Layer](#interface-layer)
10. [Database Schema](#database-schema)
11. [API Endpoints](#api-endpoints)
12. [Business Rules](#business-rules)
13. [Configuration](#configuration)
14. [Deployment](#deployment)
15. [Migration Tool](#migration-tool)
16. [Makefile Commands](#makefile-commands)
17. [Development Workflow](#development-workflow)

---

## Project Overview

This is a **Pismo Code Assessment** project implementing a RESTful API for managing customer accounts and financial transactions.
The system is built using **Go 1.25.1** and follows **Clean Architecture** principles with clear separation of concerns between domain, application, infrastructure, and interface layers.

### Key Features
- Create customer accounts with Brazilian CPF/CNPJ validation
- Retrieve account information by ID
- Create financial transactions (purchases, installment purchases, withdrawals, payments)
- Automatic amount sign handling for debit/credit operations
- Request tracing with x-trace-id headers
- Swagger/OpenAPI documentation
- Docker containerization
- Database migrations with Goose

---

## Architecture

The project follows **Clean Architecture**, organized into four main layers:

```
┌─────────────────────────────────────────────────────────┐
│                    Interface Layer                      │
│         (HTTP Handlers, Routers, Middleware)            │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   Application Layer                     │
│        (Services, DTOs, Mappers, Use Cases)             │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                    Domain Layer                         │
│      (Entities, Business Rules, Repositories)           │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 Infrastructure Layer                    │
│   (Database, Logging, Config, External Services)        │
└─────────────────────────────────────────────────────────┘
```

### Architecture Principles
- **Dependency Inversion**: Inner layers don't depend on outer layers
- **Interface-based Design**: Dependencies injected through interfaces
- **Separation of Concerns**: Each layer has a single responsibility
- **Testability**: Easy to mock and test each layer independently

---

## Technology Stack

### Core Framework & Language
- **Go 1.25.1**: Primary programming language
- **Gin (v1.11.0)**: HTTP web framework
- **Context**: Request lifecycle management

### Database
- **PostgreSQL 15**: Primary database
- **Bun (v1.2.16)**: SQL ORM and query builder
- **lib/pq (v1.10.9)**: PostgreSQL driver
- **Goose (v3.26.0)**: Database migration tool

### Configuration & CLI
- **Viper (v1.21.0)**: Configuration management
- **Cobra (v1.10.2)**: CLI framework for migration tool

### Documentation
- **Swag (v1.16.6)**: Swagger documentation generation
- **gin-swagger (v1.6.1)**: Swagger UI integration

### Validation & Utilities
- **brdoc (v1.1.2)**: Brazilian document (CPF/CNPJ) validation
- **UUID (v1.6.0)**: Unique identifier generation

### Testing
- **testify (v1.11.1)**: Testing toolkit with assertions and mocks

### Containerization
- **Docker**: Application containerization
- **Docker Compose**: Multi-container orchestration

---

## Project Structure

```
pismo-code-assessment/
├── cmd/                                  # Application entry points
│   ├── api/
│   │   └── main.go                       # Main API server
│   └── migrate/
│       └── migration_tool.go             # Database migration CLI
│
├── internal/                             # Private application code
│   ├── core/                             # Core abstractions
│   │   ├── adapter/                      # Connection adapters
│   │   ├── config/                       # Configuration interfaces
│   │   ├── contextkeys/                  # Context key constants
│   │   ├── contextutils/                 # Context utilities
│   │   ├── factory/                      # Factory interfaces
│   │   └── logger/                       # Logger interfaces
│   │
│   ├── domains/                          # Domain layer (business logic)
│   │   ├── account/
│   │   │   ├── entity.go                 # Account entity & validation
│   │   │   ├── service.go                # Account service interface
│   │   │   ├── repository.go             # Account repository interface
│   │   │   ├── mocks.go                  # Test mocks
│   │   │   └── entity_test.go            # Entity tests
│   │   └── transaction/
│   │       ├── entity.go                 # Transaction entity
│   │       ├── service.go                # Transaction service interface
│   │       ├── repository.go             # Transaction repository interface
│   │       ├── mocks.go                  # Test mocks
│   │       └── entity_test.go            # Entity tests
│   │
│   └── infra/                            # Infrastructure implementations
│       ├── config/
│       │   └── viper_config.go           # Viper config implementation
│       ├── logger/
│       │   └── slog_logger.go            # Structured logging
│       ├── factory/
│       │   └── app_factory.go            # Dependency injection factory
│       └── database/
│           ├── connection/
│           │   └── postgres_connection.go # PostgreSQL connection
│           ├── migrations/                # SQL migration files
│           │   ├── 01_create_tables.sql
│           │   └── 02_insert_operation_type.sql
│           ├── model/                     # Database models
│           │   ├── account_model.go
│           │   ├── transaction_model.go
│           │   └── operation_type_model.go
│           ├── mapper/                    # Entity-Model mappers
│           │   ├── account_mapper.go
│           │   ├── transaction_mapper.go
│           │   └── operation_type_mapper.go
│           └── repository/                # Repository implementations
│               ├── account_postgres_repository.go
│               └── transaction_postgres_repository.go
│
├── application/                           # Application layer (use cases)
│   ├── account/
│   │   ├── dto/                           # Data Transfer Objects
│   │   │   └── dto.go
│   │   ├── mapper/                        # DTO-Entity mappers
│   │   │   └── account_mapper.go
│   │   └── service/
│   │       └── account_service.go         # Account use cases
│   └── transaction/
│       ├── dto/
│       │   └── dto.go
│       ├── mapper/
│       │   └── transaction_mapper.go
│       └── service/
│           └── transaction_service.go     # Transaction use cases
│
├── interfaces/                            # Interface layer (API)
│   └── http/
│       ├── handler/                       # HTTP handlers
│       │   ├── account_handler.go
│       │   └── transaction_handler.go
│       ├── middleware/                    # HTTP middleware
│       │   ├── trace_middleware.go        # Request tracing
│       │   └── logger_middleware.go       # Request logging
│       └── router/                        # Route configuration
│           ├── router.go
│           └── router_factory.go
│
├── docs/                                  # Generated Swagger docs
│   ├── swagger.yaml
│   ├── swagger.json
│   └── docs.go
│
├── config.yaml                            # Application configuration
├── sample.config.yaml                     # Sample configuration
├── docker-compose.yaml                    # Docker Compose setup
├── Dockerfile                             # Application container
├── go.mod                                 # Go module dependencies
└── README.md                              # Project README
```

---

## Domain Layer

The domain layer contains the core business logic and is independent of any external frameworks or libraries.

### Account Domain

**Location**: `internal/domains/account/`

#### Account Entity (`entity.go`)
```go
type Account struct {
    AccountID      int64  // Unique identifier
    DocumentNumber string // Brazilian CPF or CNPJ
}
```

**Key Functions**:
- `IsValidDocumentNumber(documentNumber string) error`: Validates Brazilian CPF or CNPJ using the `brdoc` library
- `SanitizeDocumentNumber(documentNumber string) string`: Removes non-digit characters from document numbers

#### Account Service Interface (`service.go`)
```go
type Service interface {
    FindByID(ctx context.Context, request dto.FindAccountByIdRequest) (*dto.FindAccountByIdResponse, error)
    Create(ctx context.Context, response dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
}
```

**Custom Errors**:
- `AccountServiceInvalidParametersError`: Invalid input parameters
- `AccountServiceNotFoundError`: Account not found
- `AccountServiceAlreadyExistsForDocumentNumberError`: Duplicate document number

#### Account Repository Interface (`repository.go`)
```go
type AccountRepository interface {
    FindByID(ctx context.Context, accountID int64) (*Account, error)
    FindByDocumentNumber(ctx context.Context, documentNumber string) (*Account, error)
    Save(ctx context.Context, newAccount *Account) (*Account, error)
}
```

**Custom Errors**:
- `AccountRepositoryInvalidParametersError`: Invalid parameters
- `AccountRepositoryNotFoundError`: Account not found

---

### Transaction Domain

**Location**: `internal/domains/transaction/`

#### Transaction Entity (`entity.go`)
```go
type Transaction struct {
    TransactionID   int64     // Unique identifier
    AccountID       int64     // Account reference
    OperationTypeID int       // Type of operation
    Amount          float64   // Transaction amount
    EventDate       time.Time // Transaction timestamp
}

type OperationType struct {
    OperationTypeID int64
    Description     string
}
```

**Operation Type Constants**:
```go
const (
    Purchase            = 1  // Normal purchase (debit)
    InstallmentPurchase = 2  // Installment purchase (debit)
    Withdrawal          = 3  // Cash withdrawal (debit)
    Payment             = 4  // Payment (credit)
)
```

#### Transaction Service Interface (`service.go`)
```go
type Service interface {
    Create(ctx context.Context, input dto.CreateTransactionRequest) (*dto.CreateTransactionResponse, error)
}
```

**Custom Errors**:
- `TransactionServiceInvalidParametersError`: Invalid parameter
- `TransactionServiceInvalidOperationTypeError`: Invalid operation type
- `TransactionServiceInvalidAmountNegativeError`: Amount must be positive
- `TransactionServiceInvalidAccountIDError`: Invalid account ID

#### Transaction Repository Interface (`repository.go`)
```go
type TransactionRepository interface {
    FindOperationTypeByID(ctx context.Context, operationTypeID int) (*OperationType, error)
    Save(ctx context.Context, newTransaction *Transaction) (*Transaction, error)
}
```

**Custom Errors**:
- `OperationTypeRepositoryNotFoundError`: Operation type not found
- `TransactionRepositoryInvalidParametersError`: Invalid parameters

---

## Application Layer

The application layer orchestrates use cases by coordinating domain entities and repositories.

### Account Service

**Location**: `application/account/service/account_service.go`

**Implementation**: `AccountService`

**Key Methods**:

1. **FindByID** (`account_service.go:24`)
   - Validates account ID (must be > 0)
   - Queries repository for account
   - Maps entity to response DTO
   - Returns error if not found

2. **Create** (`account_service.go:40`)
   - Validates document number using Brazilian CPF/CNPJ validation
   - Checks if account already exists for document number
   - Creates new account if validation passes
   - Returns created account or error

**Dependencies**:
- `AccountRepository`: For data persistence
- `Logger`: For structured logging

---

### Transaction Service

**Location**: `application/transaction/service/transaction_service.go`

**Implementation**: `TransactionService`

**Key Methods**:

1. **Create** (`transaction_service.go:29`)
   - Validates request parameters
   - Verifies account exists
   - Reverses amount sign for debit operations
   - Saves transaction to repository
   - Returns transaction with original sign for display

2. **validateRequestParameters** (`transaction_service.go:64`)
   - Checks account ID > 0
   - Validates amount > 0
   - Verifies operation type exists

3. **reverseAmountSign** (`transaction_service.go:78`)
   - Converts positive amounts to negative for debits (operations 1-3)
   - Keeps positive for credits (operation 4)
   - Implements business rule: debits stored as negative values

**Business Logic**:
- **Debit Operations** (Purchase, Installment Purchase, Withdrawal): Store amounts as negative values
- **Credit Operations** (Payment): Store amounts as positive values
- User always sends positive amounts; service handles sign conversion

**Dependencies**:
- `AccountRepository`: To verify account exists
- `TransactionRepository`: For transaction persistence
- `Logger`: For structured logging

---

### DTOs (Data Transfer Objects)

#### Account DTOs (`application/account/dto/dto.go`)
```go
type CreateAccountRequest struct {
    DocumentNumber string `json:"document_number" binding:"required"`
}

type CreateAccountResponse struct {
    AccountID      int64  `json:"account_id"`
    DocumentNumber string `json:"document_number"`
}

type FindAccountByIdRequest struct {
    AccountID int64 `uri:"account_id" binding:"required,gt=0"`
}

type FindAccountByIdResponse struct {
    AccountID      int64  `json:"account_id"`
    DocumentNumber string `json:"document_number"`
}
```

#### Transaction DTOs (`application/transaction/dto/dto.go`)
```go
type CreateTransactionRequest struct {
    AccountID       int64   `json:"account_id"`
    OperationTypeID int     `json:"operation_type_id"`
    Amount          float64 `json:"amount"`
}

type CreateTransactionResponse struct {
    TransactionID   int64   `json:"transaction_id"`
    AccountID       int64   `json:"account_id"`
    OperationTypeID int     `json:"operation_type_id"`
    Amount          float64 `json:"amount"`
}
```

---

## Infrastructure Layer

### Database Connection

**Location**: `internal/infra/database/connection/postgres_connection.go`

**Implementation**: `PostgresConnection`
- Uses Bun ORM with PostgreSQL dialect
- Connection string from configuration
- Returns `ConnectionData` with Bun DB instance

---

### Repositories

#### Account Repository (`account_postgres_repository.go`)

**Methods**:

1. **FindByID** (`account_postgres_repository.go:29`)
   - Queries accounts table by account_id
   - Returns `AccountRepositoryNotFoundError` if not found
   - Maps database model to domain entity

2. **FindByDocumentNumber** (`account_postgres_repository.go:41`)
   - Queries accounts table by document_number
   - Returns error if not found
   - Uses unique constraint on document_number

3. **Save** (`account_postgres_repository.go:53`)
   - Inserts new account in transaction
   - Uses `Returning("*")` to get generated ID
   - Commits transaction on success
   - Rolls back on error

#### Transaction Repository (`transaction_postgres_repository.go`)

**Methods**:

1. **Save** (`transaction_postgres_repository.go:29`)
   - Inserts new transaction in database transaction
   - Returns generated transaction ID
   - Maps model to entity

2. **FindOperationTypeByID** (`transaction_postgres_repository.go:48`)
   - Queries operation_types table
   - Validates operation type exists
   - Returns operation type entity

---

### Configuration

**Location**: `internal/infra/config/viper_config.go`

**Implementation**: Uses Viper for configuration management

**Configuration Structure**:
```yaml
app:
  env: "development"
  address: ":8080"
  loglevel: "debug"

database:
  dsn: "postgres://user:password@host:port/database?sslmode=disable"
```

**Config Paths**:
- Checks environment variable `CONFIG_PATH`
- Falls back to current directory
- Looks for `config.yaml`

---

### Logging

**Location**: `internal/infra/logger/slog_logger.go`

**Implementation**: Uses Go's structured logging (slog)

**Features**:
- Structured logging with key-value pairs
- Log level configuration (debug, info, warn, error)
- Component name tracking
- Request tracing support

---

### Factory Pattern

**Location**: `internal/infra/factory/app_factory.go`

**Implementation**: `AppFactory`

**Purpose**: Dependency injection container

**Methods**:
- `AccountRepository()`: Creates account repository
- `TransactionRepository()`: Creates transaction repository
- `AccountService()`: Creates account service with dependencies
- `TransactionService()`: Creates transaction service with dependencies
- `AccountHandler()`: Creates account HTTP handler
- `TransactionHandler()`: Creates transaction HTTP handler

**Benefits**:
- Centralized dependency management
- Easy to swap implementations
- Supports testing with mocks

---

## Interface Layer

### HTTP Handlers

#### Account Handler

**Location**: `interfaces/http/handler/account_handler.go`

**Endpoints**:

1. **CreateAccount** (`account_handler.go:34`)
   - Method: POST
   - Path: `/accounts`
   - Binds JSON request to DTO
   - Calls account service
   - Returns 201 Created or 400 Bad Request

2. **GetAccountByID** (`account_handler.go:57`)
   - Method: GET
   - Path: `/accounts/:account_id`
   - Parses account ID from URL parameter
   - Calls account service
   - Returns 200 OK, 404 Not Found, or 500 Internal Server Error

#### Transaction Handler

**Location**: `interfaces/http/handler/transaction_handler.go`

**Endpoints**:

1. **CreateTransaction** (`transaction_handler.go:33`)
   - Method: POST
   - Path: `/transactions`
   - Binds JSON request to DTO
   - Calls transaction service
   - Returns 201 Created or 400 Bad Request

---

### Middleware

#### Trace Middleware

**Location**: `interfaces/http/middleware/trace_middleware.go`

**Purpose**: Request tracing with x-trace-id header

**Behavior**:
- Checks for existing `x-trace-id` header
- Generates new UUID if not present
- Adds to response headers
- Stores in request context for logging

#### Logger Middleware

**Location**: `interfaces/http/middleware/logger_middleware.go`

**Purpose**: Request/response logging

**Features**:
- Logs request method, path, status code
- Includes trace ID
- Tracks request duration

---

### Router

**Location**: `interfaces/http/router/router_factory.go` and `router.go`

**Configuration**:
- Sets up Gin router
- Applies middleware (tracing, logging)
- Registers account routes
- Registers transaction routes
- Serves Swagger documentation at `/swagger/*`

**Route Structure**:
```
POST   /accounts
GET    /accounts/:account_id
POST   /transactions
GET    /swagger/*  (Swagger UI)
```

---

## Database Schema

### Tables

#### accounts
```sql
CREATE TABLE accounts (
    account_id      BIGSERIAL PRIMARY KEY,
    document_number VARCHAR NOT NULL UNIQUE
);
```

**Indexes**:
- Primary key on `account_id`
- Unique constraint on `document_number`

---

#### operation_types
```sql
CREATE TABLE operation_types (
    operation_type_id BIGINT PRIMARY KEY,
    description       VARCHAR NOT NULL
);
```

**Seed Data**:
```sql
INSERT INTO operation_types VALUES (1, 'PURCHASE');
INSERT INTO operation_types VALUES (2, 'INSTALLMENT PURCHASE');
INSERT INTO operation_types VALUES (3, 'WITHDRAWAL');
INSERT INTO operation_types VALUES (4, 'PAYMENT');
```

---

#### transactions
```sql
CREATE TABLE transactions (
    transaction_id BIGSERIAL PRIMARY KEY,
    account_id     BIGINT NOT NULL REFERENCES accounts(account_id),
    operation_type BIGINT NOT NULL REFERENCES operation_types(operation_type_id),
    amount         DOUBLE PRECISION NOT NULL,
    event_date     TIMESTAMP WITH TIME ZONE NOT NULL
);
```

**Foreign Keys**:
- `account_id` → `accounts.account_id`
- `operation_type` → `operation_types.operation_type_id`

---

## API Endpoints

### Swagger Documentation

**URL**: `http://localhost:8080/swagger/index.html`

**Metadata**:
- Title: Pismo Code Assessment API
- Version: 1.0
- Description: Customer Account & Transactions
- Author: Fábio Sartori
- License: MIT

---

### Create Account

**Endpoint**: `POST /accounts`

**Request Body**:
```json
{
  "document_number": "12345678900"
}
```

**Response (201 Created)**:
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

**Errors**:
- 400 Bad Request: Invalid document number or duplicate

---

### Get Account by ID

**Endpoint**: `GET /accounts/{id}`

**Response (200 OK)**:
```json
{
  "account_id": 1,
  "document_number": "12345678900"
}
```

**Errors**:
- 404 Not Found: Account doesn't exist

---

### Create Transaction

**Endpoint**: `POST /transactions`

**Request Body**:
```json
{
  "account_id": 1,
  "operation_type_id": 1,
  "amount": 123.45
}
```

**Response (201 Created)**:
```json
{
  "transaction_id": 1,
  "account_id": 1,
  "operation_type_id": 1,
  "amount": 123.45
}
```

**Errors**:
- 400 Bad Request: Invalid parameters, account not found, or invalid operation type

---

## Business Rules

### Document Validation
- Only Brazilian CPF (11 digits) or CNPJ (14 digits) are accepted
- Document numbers must be valid according to Brazilian checksum algorithms
- Non-digit characters are automatically stripped
- Document numbers must be unique across accounts

### Transaction Amount Handling
- Users always send positive amounts
- System automatically converts to negative for debit operations:
  - Operation 1 (Purchase): Stored as negative
  - Operation 2 (Installment Purchase): Stored as negative
  - Operation 3 (Withdrawal): Stored as negative
  - Operation 4 (Payment): Stored as positive
- Amounts returned in responses maintain display sign (positive)

### Operation Types
1. **Purchase**: Regular purchase transaction (debit)
2. **Installment Purchase**: Purchase paid in installments (debit)
3. **Withdrawal**: Cash withdrawal (debit)
4. **Payment**: Payment/deposit (credit)

### Validation Rules
- Account ID must be > 0
- Amount must be > 0 (positive)
- Operation type must exist in database
- Account must exist before creating transactions
- Document number must not already be registered

---

## Configuration

### Application Configuration

**File**: `config.yaml`

**Structure**:
```yaml
app:
  env: "development"          # Environment (development/production)
  address: ":8080"            # Server bind address
  loglevel: "debug"           # Log level (debug/info/warn/error)

database:
  dsn: "postgres://pismo:123456@localhost:5432/pismo_db?sslmode=disable"
```

**Environment Variable**:
- `CONFIG_PATH`: Custom path to configuration file

---

## Deployment

### Docker Compose Setup

**File**: `docker-compose.yaml`

**Services**:

1. **postgres**:
   - Image: `postgres:15-alpine`
   - Port: 5432
   - Credentials:
     - User: `pismo`
     - Password: `123456`
     - Database: `pismo_db`
   - Healthcheck: Waits for PostgreSQL to be ready

2. **pismo-api**:
   - Built from Dockerfile
   - Port: 8080
   - Depends on postgres service
   - Mounts config.yaml
   - Auto-restarts on failure

**Commands**:
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f pismo-api

# Stop all services
docker-compose down

# Rebuild and start
docker-compose up --build
```

---

### Dockerfile

**Multi-stage Build**:

**Stage 1 - Builder**:
```dockerfile
FROM golang:1.25.6-alpine3.23 as builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api
```

**Stage 2 - Runtime**:
```dockerfile
FROM alpine:3.23.2
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

**Benefits**:
- Small final image (Alpine Linux)
- Excludes build dependencies from runtime
- Fast builds with layer caching

---

## Migration Tool

**Location**: `cmd/migrate/migration_tool.go`

**Purpose**: CLI tool for database migrations using Goose

### Commands

#### Run Migrations
```bash
go run cmd/migrate/migration_tool.go up
```
- Applies all pending migrations
- Creates tables and seeds data

#### Rollback All Migrations
```bash
go run cmd/migrate/migration_tool.go down
```
- Prompts for confirmation
- Destroys all database tables
- Use with caution

#### Check Migration Status
```bash
go run cmd/migrate/migration_tool.go status
```
- Shows applied and pending migrations
- Displays migration versions

### Migration Files

**Location**: `internal/infra/database/migrations/`

1. **01_create_tables.sql**: Creates accounts, operation_types, and transactions tables
2. **02_insert_operation_type.sql**: Seeds operation types (1-4)

**Format**: Goose SQL migrations with `+goose up` and `+goose down` sections

---

i## Makefile Commands

**Location**: `Makefile`

**Purpose**: Provides convenient shortcuts for common development tasks

The Makefile contains several targets to simplify development, testing, deployment, and database management workflows.

### Application Commands

#### Run API Server
```bash
make run
```
- Starts the API server locally
- Executes `go run cmd/api/main.go`
- Requires database to be running and migrated

#### Install Dependencies
```bash
make install
```
- Tidies up Go module dependencies
- Downloads all required packages
- Creates vendor directory with dependencies
- Executes:
  - `go mod tidy`: Removes unused dependencies
  - `go mod vendor`: Copies dependencies to vendor folder

---

### Database Migration Commands

#### Apply All Migrations
```bash
make db-migration-up
```
- Runs all pending database migrations
- Creates tables and seeds initial data
- Executes: `go run cmd/migrate/migration_tool.go up`

#### Rollback All Migrations
```bash
make db-migration-down
```
- Prompts for confirmation before destroying all tables
- Rolls back all migrations to initial state
- Executes: `go run cmd/migrate/migration_tool.go down`
- **Use with caution**: This will delete all data

#### Check Migration Status
```bash
make db-migration-status
```
- Displays current migration status
- Shows which migrations are applied and pending
- Executes: `go run cmd/migrate/migration_tool.go status`

---

### Docker-Based Migration Commands

These commands run migrations inside the Docker container, useful when the API is running in Docker.

#### Apply Migrations Inside Docker Container
```bash
make db-migration-docker-up
```
- Runs all pending database migrations inside the running Docker container
- Requires the container to be running (`make docker-up`)
- Executes: `docker exec -it pismo-api /app/app_migration_tool up`
- Used automatically by `make docker-provision`

#### Rollback Migrations Inside Docker Container
```bash
make db-migration-docker-down
```
- Rolls back all migrations inside the running Docker container
- Prompts for confirmation before destroying all tables
- Executes: `docker exec -it pismo-api /app/app_migration_tool down`
- **Use with caution**: This will delete all data

#### Check Migration Status Inside Docker Container
```bash
make db-migration-docker-status
```
- Displays current migration status from inside the Docker container
- Shows which migrations are applied and pending
- Executes: `docker exec -it pismo-api /app/app_migration_tool status`

---

### Testing Commands

#### Run Unit Tests
```bash
make test-unit
```
- Runs all unit tests in the project
- Displays verbose output
- Executes: `go test -v ./...`
- Tests located in:
  - `internal/domains/account/entity_test.go`
  - `internal/domains/transaction/entity_test.go`
  - `internal/infra/config/viper_config_test.go`

---

### Docker Commands

#### Complete Docker Provisioning (Recommended)
```bash
make docker-provision
```
- **This is the recommended way to start the application**
- Builds and starts all Docker containers
- Automatically runs all database migrations inside the container
- Single command to get everything running
- Executes:
  1. `make docker-up`: Starts all containers
  2. `make db-migration-docker-up`: Runs migrations in container
- Perfect for first-time setup or quick deployments

#### Build and Start All Services
```bash
make docker-up
```
- Builds Docker images
- Starts all services defined in docker-compose.yaml
- Runs in detached mode (-d)
- Includes:
  - PostgreSQL database
  - Pismo API application
- Executes: `docker compose up -d --build`
- **Note**: This does NOT run migrations automatically. Use `docker-provision` instead or run migrations manually with `make db-migration-docker-up`

#### Destroy All Containers
```bash
make docker-destroy
```
- Stops and removes all containers
- Removes orphaned containers
- Deletes all volumes (including database data)
- Executes: `docker compose down --remove-orphans --volumes`
- **Warning**: This will delete all database data

#### Start PostgreSQL Only
```bash
make docker-up-postgres
```
- Starts only the PostgreSQL database container
- Useful for local development without containerizing the API
- Executes: `docker compose up -d postgres`

#### Stop PostgreSQL
```bash
make docker-down-db
```
- Stops and removes the PostgreSQL container
- Executes: `docker compose down -d postgres`

---

### Documentation Commands

#### Generate Swagger Documentation
```bash
make swagger
```
- Generates Swagger/OpenAPI documentation
- Parses annotations in Go code
- Creates documentation files in `docs/` directory:
  - `docs/swagger.json`
  - `docs/swagger.yaml`
  - `docs/docs.go`
- Executes: `swag init -g cmd/api/main.go -o docs`
- Run this after adding or modifying API endpoint annotations

#### Generate Configuration File
```bash
make config-file
```
- Generates a new `config.yaml` file with default values
- Creates the file in a temporary folder
- Useful for creating template configuration files
- Executes: `go run cmd/config/main.go`
- The generated config includes default values for:
  - Server port and address
  - Database connection settings
  - Logger configuration
  - Environment settings

---

### Makefile Variables

The Makefile defines the following variables:

```makefile
MIGRATION_TOOL = "cmd/migrate/migration_tool.go"      # Path to migration tool source
MIGRATION_TOOL_BIN = "/app/app_migration_tool"        # Migration tool binary path in container
IMAGE_NAME = pismo-api                                 # Docker image name
CONTAINER_NAME = pismo-api                             # Docker container name
```

---

### Common Development Workflows Using Makefile

#### Quick Start (Recommended)
```bash
# Single command to run everything with Docker
make docker-provision
```
This is the fastest way to get the application running. It handles everything:
- Builds Docker images
- Starts all containers (API + PostgreSQL)
- Runs all database migrations
- Application ready at http://localhost:8080

#### First Time Setup (Local Development)
```bash
# 1. Install dependencies
make install

# 2. Start database
make docker-up-postgres

# 3. Run migrations
make db-migration-up

# 4. Start API server
make run
```

#### Full Docker Deployment (Manual Steps)
```bash
# 1. Build and start all containers
make docker-up

# 2. Run migrations inside container
make db-migration-docker-up

# 3. Verify it's running
docker ps

# 4. View logs
docker compose logs -f pismo-api
```
**Note**: Using `make docker-provision` combines steps 1 and 2 automatically.

#### Testing Workflow
```bash
# Run tests
make test-unit

# If tests pass, generate updated documentation
make swagger
```

#### Database Reset
```bash
# Destroy and recreate database
make db-migration-down  # Rollback all
make db-migration-up    # Apply all again
```

#### Clean Shutdown
```bash
# Stop and clean up everything
make docker-destroy
```

---

## Development Workflow

### Local Development

1. **Install Dependencies**:
   ```bash
   go mod download
   ```

2. **Start Database**:
   ```bash
   docker-compose up -d postgres
   ```

3. **Run Migrations**:
   ```bash
   go run cmd/migrate/migration_tool.go up
   ```

4. **Run API Server**:
   ```bash
   go run cmd/api/main.go
   ```

5. **Access Swagger**:
   ```
   http://localhost:8080/swagger/index.html
   ```

### Testing

**Run Tests**:
```bash
go test ./...
```

**Run Tests with Coverage**:
```bash
go test -cover ./...
```

**Test Files**:
- `internal/domains/account/entity_test.go`
- `internal/domains/transaction/entity_test.go`
- `internal/infra/config/viper_config_test.go`

---

## Key Design Patterns

### 1. Clean Architecture
- Separation of domain, application, infrastructure, and interface layers
- Dependency inversion: inner layers don't depend on outer layers

### 2. Repository Pattern
- Abstracts data access logic
- Domain defines interfaces, infrastructure implements

### 3. Factory Pattern
- `AppFactory`: Centralized dependency injection
- `RouterFactory`: Router configuration
- Simplifies object creation and wiring

### 4. DTO Pattern
- Separate DTOs for API layer
- Decouples API contracts from domain entities
- Mappers convert between DTOs and entities

### 5. Middleware Pattern
- Composable request processing (tracing, logging)
- Cross-cutting concerns separated from business logic

### 6. Interface Segregation
- Small, focused interfaces
- Easy to mock for testing
- Services depend on interfaces, not concrete implementations

---

## Error Handling Strategy

### Domain Errors
- Defined as package-level variables
- Descriptive error names (e.g., `AccountServiceNotFoundError`)
- Propagated up through layers

### HTTP Error Mapping
- Domain errors mapped to HTTP status codes
- Consistent error response format: `{"error": "message"}`
- Status codes:
  - 200: Success
  - 201: Created
  - 400: Bad Request (validation errors)
  - 404: Not Found
  - 500: Internal Server Error

### Transaction Safety
- Database transactions used for writes
- Deferred rollback ensures cleanup
- Explicit commit on success

---

## Logging Strategy

### Structured Logging
- Uses Go's `slog` package
- Key-value pairs for context
- Component name tracking

### Log Levels
- **Debug**: Detailed operation logs (parameters, queries)
- **Info**: General application flow
- **Warn**: Warning conditions
- **Error**: Error conditions

### Request Tracing
- x-trace-id header in all requests
- Trace ID propagated through context
- Enables request tracking across services

---

## Security Considerations

### Input Validation
- JSON binding with required fields
- Document number validation (CPF/CNPJ)
- Positive amount validation
- Operation type validation

### SQL Injection Prevention
- Uses parameterized queries (Bun ORM)
- No raw SQL string concatenation

### Database Credentials
- Configuration file-based (not hardcoded)
- Should use environment variables in production
- Docker secrets recommended for production

---

## Performance Considerations

### Database
- Indexes on primary keys
- Unique constraint on document_number
- Transaction pooling via Bun

### Application
- Efficient JSON serialization (Gin)
- Context-based request cancellation
- Connection pooling for database

### Docker
- Multi-stage builds for smaller images
- Alpine Linux for minimal footprint
- Healthchecks for reliability

---

## Conclusion

This project demonstrates a well-structured, maintainable Go application following Clean Architecture principles. It showcases:

- **Clear separation of concerns** across layers
- **Interface-based design** for flexibility and testability
- **Domain-driven design** with rich domain models
- **Infrastructure independence** in business logic
- **Production-ready practices**: Docker, migrations, logging, documentation
- **Brazilian market specificity**: CPF/CNPJ validation
- **Sound business rules**: Automatic debit/credit handling

The codebase is designed for scalability, maintainability, and extensibility, making it easy to add new features or swap implementations without major refactoring.

---

**Author**: Fábio Sartori
**Repository**: https://github.com/kiosanim
**License**: MIT
