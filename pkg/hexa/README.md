# Hexagonal Architecture Example

This is a demonstration of Hexagonal Architecture (Ports & Adapters) using Go.

## Concepts

The core principle is that **dependencies point inward**. The application logic (Domain) knows nothing about the outside world.

- **Domain (`domain/`)**: The core business logic.
  - **Entities**: Business objects (e.g., `User`). They have no external dependencies (not even for UUID generation).
  - **Ports**: Interfaces describing what the application can do (Primary) and what it needs (Secondary).
- **Use Cases (`domain/usecases/`)**: Concrete implementations of Primary Ports, orchestrating the business rules.
- **Adapters (`adapters/`)**: Concrete implementations of Secondary Ports (e.g., Postgres, In-Memory DB, Console Mailer, UUID Generation directly in use cases).
- **Hosts (`hosts/`)**: The driving adapters that trigger the application. We demonstrate two independent entrypoints:
  - **HTTP API (`hosts/http/`)**: Listens for HTTP requests, parses them into DTOs, and calls the Use Case.
  - **CLI (`hosts/cli/`)**: Takes terminal arguments and calls the exact same Use Case.
- **Main (`main.go`)**: The Composition Root that wires everything together.

## Where to add things?

- **New business logic / features**:
  1. Define a new method in the Primary Port interface in `domain/ports.go`.
  2. Implement the method in the Use Case (`domain/usecases/`).
  3. Write tests for the Use Case.
- **New data source / external service / utility**:
  2. Implement the interface in a new adapter struct in `adapters/`.
- **New API endpoint**:
  1. Add Request/Response DTOs in the HTTP host (`hosts/http/user_handler.go`) to prevent leaking the Domain entities.
  2. Add a route and handler function, map the Request DTO to Domain inputs, call the Primary Port, and map the Domain output to a Response DTO.

## Running the Application

### 1. HTTP Server
```bash
go run main.go
```
*(Starts an HTTP server on `:8080`)*

### 2. CLI Tool
```bash
go run hosts/cli/main.go --email="cli@test.com" --name="CLI User"
```
*(Executes the same registration use case directly from the terminal!)*

## Running Tests

This project includes both unit tests (using mocks and in-memory adapters) and integration tests (using `testcontainers-go` for PostgreSQL).

```bash
go test ./...
```
*(Requires Docker to be running for the integration tests)*
