package domain

// Primary port — any caller (HTTP, CLI, test) uses this interface.
type UserUseCase interface {
    RegisterUser(email, name string) (User, error)
    GetUser(id string) (User, error)
    ListUsers() ([]User, error)
    UpdateUser(id, email, name string) (User, error)
    DeleteUser(id string) error
}

// Secondary port — storage. The domain defines the contract; adapters fulfil it.
type UserRepository interface {
    Save(user User) error
    ExistsByEmail(email string) (bool, error)
    FindByID(id string) (User, error)
    FindAll() ([]User, error)
    Update(user User) error
    Delete(id string) error
}

// Secondary port — notifications.
type WelcomeMailer interface {
    SendWelcome(toEmail, name string) error
}
