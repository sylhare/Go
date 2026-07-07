package adapters_test

import (
    "context"
    "database/sql"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"

    _ "github.com/lib/pq"
    "hexa/adapters"
    "hexa/domain"
)

func setupPostgres(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()

    dbName := "users_db"
    dbUser := "user"
    dbPassword := "password"

    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase(dbName),
        postgres.WithUsername(dbUser),
        postgres.WithPassword(dbPassword),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).WithStartupTimeout(5*time.Second),
        ),
    )
    require.NoError(t, err)

    connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)

    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)

    // Create table
    _, err = db.Exec(`
        CREATE TABLE users (
            id VARCHAR(255) PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            name VARCHAR(255) NOT NULL
        )
    `)
    require.NoError(t, err)

    cleanup := func() {
        db.Close()
        postgresContainer.Terminate(ctx)
    }

    return db, cleanup
}

func TestPostgresUserRepository_CRUD(t *testing.T) {
    db, cleanup := setupPostgres(t)
    defer cleanup()

    repo := adapters.NewPostgresUserRepository(db)
    user := domain.User{ID: "123", Email: "test@example.com", Name: "Test"}

    // Create
    err := repo.Save(user)
    require.NoError(t, err)

    // Exists
    exists, err := repo.ExistsByEmail("test@example.com")
    require.NoError(t, err)
    assert.True(t, exists)

    // Read
    fetched, err := repo.FindByID("123")
    require.NoError(t, err)
    assert.Equal(t, user, fetched)

    // Update
    user.Name = "Test Updated"
    err = repo.Update(user)
    require.NoError(t, err)

    fetched, _ = repo.FindByID("123")
    assert.Equal(t, "Test Updated", fetched.Name)

    // FindAll
    users, err := repo.FindAll()
    require.NoError(t, err)
    assert.Len(t, users, 1)

    // Delete
    err = repo.Delete("123")
    require.NoError(t, err)

    // Read after delete
    _, err = repo.FindByID("123")
    require.ErrorIs(t, err, domain.ErrUserNotFound)
}
