# Hexagonal Architecture Example

This is a demonstration of Hexagonal Architecture (Ports & Adapters) using Go.

## Concepts

The core principle is that **dependencies point inward**. The application logic (Domain) knows nothing about the outside world.

- **Domain (`domain/`)**: The core business logic.
  - **Entities**: Business objects (e.g., `User`).
  - **Ports**: Interfaces describing what the application can do (Primary) and what it needs (Secondary).
- **Use Cases (`domain/usecases/`)**: Concrete implementations of Primary Ports, orchestrating the business rules.
- **Adapters (`adapters/`)**: Concrete implementations of Secondary Ports (e.g., Postgres, In-Memory DB, SMTP Mailer).
- **Hosts (`hosts/`)**: The driving adapters that trigger the application (e.g., HTTP API, CLI).
- **Main (`main.go`)**: The Composition Root that wires everything together.

## Where to add things?

- **New business logic / features**:
  1. Define a new method in the Primary Port interface in `domain/ports.go`.
  2. Implement the method in the Use Case (`domain/usecases/`).
  3. Write tests for the Use Case.
- **New data source / external service**:
  1. Define a Secondary Port interface in `domain/ports.go` if it doesn't exist.
  2. Implement the interface in a new adapter struct in `adapters/`.
- **New API endpoint**:
  1. Add a route and handler function in the HTTP host (`hosts/http/user_handler.go`).
  2. Call the corresponding Primary Port method.

## Running Tests

This project includes both unit tests and integration tests using `testcontainers-go` for PostgreSQL.

```bash
go test ./...
```
*(Requires Docker to be running for the integration tests)*
