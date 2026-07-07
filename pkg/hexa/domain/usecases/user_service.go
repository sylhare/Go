package usecases

import (
    "fmt"

    "hexa/domain"
    "github.com/google/uuid"
)

// UserService implements domain.UserUseCase.
type UserService struct {
    users  domain.UserRepository
    mailer domain.WelcomeMailer
}

func NewUserService(users domain.UserRepository, mailer domain.WelcomeMailer) *UserService {
    return &UserService{users: users, mailer: mailer}
}

func (s *UserService) RegisterUser(email, name string) (domain.User, error) {
    exists, err := s.users.ExistsByEmail(email)
    if err != nil {
        return domain.User{}, fmt.Errorf("checking email: %w", err)
    }
    if exists {
        return domain.User{}, fmt.Errorf("email already registered")
    }

    user, err := domain.NewUser(uuid.NewString(), email, name)
    if err != nil {
        return domain.User{}, err
    }

    if err := s.users.Save(user); err != nil {
        return domain.User{}, fmt.Errorf("saving user: %w", err)
    }

    _ = s.mailer.SendWelcome(user.Email, user.Name) // best-effort
    return user, nil
}

func (s *UserService) GetUser(id string) (domain.User, error) {
    return s.users.FindByID(id)
}

func (s *UserService) ListUsers() ([]domain.User, error) {
    return s.users.FindAll()
}

func (s *UserService) UpdateUser(id, email, name string) (domain.User, error) {
    user, err := s.users.FindByID(id)
    if err != nil {
        return domain.User{}, err
    }

    // if email changed, verify it's not taken
    if email != "" && user.Email != email {
        exists, err := s.users.ExistsByEmail(email)
        if err != nil {
            return domain.User{}, fmt.Errorf("checking email: %w", err)
        }
        if exists {
            return domain.User{}, fmt.Errorf("email already registered")
        }
    }

    // use NewUser to validate new values
    updatedUser, err := domain.NewUser(id, email, name)
    if err != nil {
        return domain.User{}, err
    }

    if err := s.users.Update(updatedUser); err != nil {
        return domain.User{}, fmt.Errorf("updating user: %w", err)
    }

    return updatedUser, nil
}

func (s *UserService) DeleteUser(id string) error {
    _, err := s.users.FindByID(id)
    if err != nil {
        return err
    }
    return s.users.Delete(id)
}
