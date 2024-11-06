package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) (*postgres.PostgresContainer, string) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("..", "db", "schema.sql")),
		//postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get container port: %s", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, host, port.Port(), dbName)
	log.Printf("dbURL: %s", dbURL)

	return postgresContainer, dbURL
}

func TestRun(t *testing.T) {
	postgresContainer, dsn := setupTestContainer(t)

	defer func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	if err := run(dsn); err != nil {
		t.Fatalf("run() failed: %s", err)
	}
}

// Need the postgres in docker running
func TestRunOnly(t *testing.T) {
	t.Skip("Skipping TestRunOnly")
	if err := run(""); err != nil {
		t.Fatalf("run() failed: %s", err)
	}
}
