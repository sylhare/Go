package domain

import (
    "errors"
    "strings"
)

var (
    ErrUserNotFound = errors.New("user not found")
)

type User struct {
    ID    string
    Email string
    Name  string
}

// NewUser enforces business invariants at construction time.
func NewUser(id, email, name string) (User, error) {
    email = strings.TrimSpace(strings.ToLower(email))
    if email == "" {
        return User{}, errors.New("email is required")
    }
    if name == "" {
        return User{}, errors.New("name is required")
    }
    return User{ID: id, Email: email, Name: name}, nil
}
