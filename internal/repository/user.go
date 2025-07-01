package repository

import (
	"context"
	"go-user-api/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	Get(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, u *model.User) error
	Delete(ctx context.Context, id int) error
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	return r.db.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)

}

func (r *UserRepo) Get(ctx context.Context, id int) (*model.User, error) {
	row := r.db.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", id)
	var u model.User

	err := row.Scan(&u.ID, &u.Name, &u.Email)
	return &u, err
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, u.ID)

	return err
}

func (r *UserRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)

	return err
}
