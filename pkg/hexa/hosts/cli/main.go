package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "hexa/adapters"
    "hexa/domain/usecases"
)

func main() {
    email := flag.String("email", "", "Email of the user to register")
    name := flag.String("name", "", "Name of the user to register")
    flag.Parse()

    if *email == "" || *name == "" {
        fmt.Println("Usage: go run main.go -email <email> -name <name>")
        os.Exit(1)
    }

    // 1. Build infrastructure (adapters — the outer ring)
    // We can use the InMemory repository here for simplicity, proving
    // the domain works regardless of the database used!
    userRepo := adapters.NewInMemoryUserRepository()
    mailer := adapters.NewConsoleWelcomeMailer()

    // 2. Build use cases (inject adapters through interfaces)
    userService := usecases.NewUserService(userRepo, mailer)

    // 3. Drive the application!
    fmt.Printf("Registering user: %s <%s>...\n", *name, *email)
    user, err := userService.RegisterUser(*email, *name)
    if err != nil {
        log.Fatalf("Failed to register user: %v\n", err)
    }

    fmt.Printf("Successfully registered user! ID: %s\n", user.ID)
}
