package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/askerdev/pinterest.user/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

const migration = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    photo_url TEXT DEFAULT ''
);
`

type User interface {
	Create(ctx context.Context, username, password string) (*domain.User, error)
	DeleteByID(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}

type user struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) User {
	db.MustExec(migration)
	return &user{db}
}

func (r *user) Create(ctx context.Context, username, password string) (*domain.User, error) {
	const query = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *`

	u := &domain.User{}
	err := r.db.GetContext(ctx, u, query, username, password)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserAlreadyExists
	}

	return u, nil
}

func (r *user) DeleteByID(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}

	return err
}

func (r *user) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `SELECT * FROM users WHERE id = $1 FETCH FIRST 1 ROW ONLY`

	u := &domain.User{}
	err := r.db.GetContext(ctx, u, query, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrUserNotFound
	case err != nil:
		return nil, err
	default:
		return u, err
	}
}

func (r *user) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	const query = `SELECT * FROM users WHERE username = $1 FETCH FIRST 1 ROW ONLY`

	u := &domain.User{}
	err := r.db.GetContext(ctx, u, query, username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, domain.ErrUserNotFound
	case err != nil:
		return nil, err
	default:
		return u, err
	}
}
