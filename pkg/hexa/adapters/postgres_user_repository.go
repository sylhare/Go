package adapters

import (
    "database/sql"

    "hexa/domain"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(user domain.User) error {
    _, err := r.db.Exec(
        `INSERT INTO users (id, email, name) VALUES ($1, $2, $3)`,
        user.ID, user.Email, user.Name,
    )
    return err
}

func (r *PostgresUserRepository) ExistsByEmail(email string) (bool, error) {
    var count int
    err := r.db.QueryRow(`SELECT COUNT(1) FROM users WHERE email = $1`, email).Scan(&count)
    return count > 0, err
}

func (r *PostgresUserRepository) FindByID(id string) (domain.User, error) {
    var user domain.User
    err := r.db.QueryRow(`SELECT id, email, name FROM users WHERE id = $1`, id).
        Scan(&user.ID, &user.Email, &user.Name)
    if err == sql.ErrNoRows {
        return domain.User{}, domain.ErrUserNotFound
    }
    return user, err
}

func (r *PostgresUserRepository) FindAll() ([]domain.User, error) {
    rows, err := r.db.Query(`SELECT id, email, name FROM users`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []domain.User
    for rows.Next() {
        var user domain.User
        if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, rows.Err()
}

func (r *PostgresUserRepository) Update(user domain.User) error {
    res, err := r.db.Exec(
        `UPDATE users SET email = $1, name = $2 WHERE id = $3`,
        user.Email, user.Name, user.ID,
    )
    if err != nil {
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return domain.ErrUserNotFound
    }
    return nil
}

func (r *PostgresUserRepository) Delete(id string) error {
    res, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, id)
    if err != nil {
        return err
    }
    rows, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rows == 0 {
        return domain.ErrUserNotFound
    }
    return nil
}
