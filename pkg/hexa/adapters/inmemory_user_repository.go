package adapters

import "hexa/domain"

type InMemoryUserRepository struct {
    store map[string]domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{store: make(map[string]domain.User)}
}

func (r *InMemoryUserRepository) Save(user domain.User) error {
    r.store[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) ExistsByEmail(email string) (bool, error) {
    for _, u := range r.store {
        if u.Email == email {
            return true, nil
        }
    }
    return false, nil
}

func (r *InMemoryUserRepository) FindByID(id string) (domain.User, error) {
    u, ok := r.store[id]
    if !ok {
        return domain.User{}, domain.ErrUserNotFound
    }
    return u, nil
}

func (r *InMemoryUserRepository) FindAll() ([]domain.User, error) {
    users := make([]domain.User, 0, len(r.store))
    for _, u := range r.store {
        users = append(users, u)
    }
    return users, nil
}

func (r *InMemoryUserRepository) Update(user domain.User) error {
    if _, ok := r.store[user.ID]; !ok {
        return domain.ErrUserNotFound
    }
    r.store[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
    if _, ok := r.store[id]; !ok {
        return domain.ErrUserNotFound
    }
    delete(r.store, id)
    return nil
}
