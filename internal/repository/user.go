package repository

import (
	"context"
	"fmt"
	"go-user-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id int) error
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", u.Name, u.Email, u.Password).Scan(&u.ID)

}

func (r *UserRepo) Get(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id)
	var u model.User

	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, err
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, name, email FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name, email, password FROM users WHERE email = $1", email)
	var u model.User

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	return &u, err
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, u.ID)

	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	res, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no users found with id: %d", id)
	}

	return err
}
