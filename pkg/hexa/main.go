package main

import (
    "database/sql"
    "net/http"
    "log"

    "hexa/adapters"
    "hexa/domain/usecases"
    httphost "hexa/hosts/http"

    _ "github.com/lib/pq"
)

func main() {
    // 1. Build infrastructure (adapters — the outer ring)
    db, err := sql.Open("postgres", "postgres://localhost/myapp?sslmode=disable")
    if err != nil {
        log.Fatalf("failed to connect to db: %v", err)
    }
    userRepo   := adapters.NewPostgresUserRepository(db)
    mailer     := adapters.NewSmtpWelcomeMailer("smtp.example.com:587")

    // 2. Build use cases (inject adapters through interfaces)
    userService := usecases.NewUserService(userRepo, mailer)

    // 3. Build hosts (inject use cases through interfaces)
    handler := httphost.NewUserHandler(userService)

    // 4. Setup Routing and start server
    mux := http.NewServeMux()
    handler.RegisterRoutes(mux)

    log.Println("Server listening on :8080")
    if err := http.ListenAndServe(":8080", mux); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
