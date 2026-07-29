package usecases_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "hexa/adapters"
    "hexa/domain/usecases"
)

type noopMailer struct{}

func (m *noopMailer) SendWelcome(toEmail, name string) error {
    return nil
}

func TestUserService_RegisterUser(t *testing.T) {
    repo   := adapters.NewInMemoryUserRepository()
    mailer := &noopMailer{}
    uc     := usecases.NewUserService(repo, mailer)

    user, err := uc.RegisterUser("alice@example.com", "Alice")
    require.NoError(t, err)
    assert.Equal(t, "alice@example.com", user.Email)
    assert.NotEmpty(t, user.ID)

    _, err = uc.RegisterUser("alice@example.com", "Alice Again")
    require.Error(t, err)
    assert.Contains(t, err.Error(), "email already registered")
}

func TestUserService_CRUD(t *testing.T) {
    repo   := adapters.NewInMemoryUserRepository()
    mailer := &noopMailer{}
    uc     := usecases.NewUserService(repo, mailer)

    // Create
    user, err := uc.RegisterUser("bob@example.com", "Bob")
    require.NoError(t, err)

    // Read
    fetched, err := uc.GetUser(user.ID)
    require.NoError(t, err)
    assert.Equal(t, user.ID, fetched.ID)

    // Update
    updated, err := uc.UpdateUser(user.ID, "bob.new@example.com", "Bob New")
    require.NoError(t, err)
    assert.Equal(t, "bob.new@example.com", updated.Email)

    // List
    users, err := uc.ListUsers()
    require.NoError(t, err)
    assert.Len(t, users, 1)

    // Delete
    err = uc.DeleteUser(user.ID)
    require.NoError(t, err)

    // Read after delete
    _, err = uc.GetUser(user.ID)
    require.Error(t, err)
}
