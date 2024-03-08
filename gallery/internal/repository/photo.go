package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/askerdev/pinterest.gallery/internal/domain"
	"github.com/jmoiron/sqlx"
)

const photoTable = `
CREATE TABLE IF NOT EXISTS photos (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id UUID NOT NULL
);
`

type Photo interface {
	Create(ctx context.Context, userID, title, url string) error
	DeleteByID(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]*domain.Photo, error)
	GetByID(ctx context.Context, id string) (*domain.Photo, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*domain.Photo, error)
}

type photo struct {
	db *sqlx.DB
}

func NewPhoto(db *sqlx.DB) Photo {
	db.Exec(photoTable)
	return &photo{db}
}

func (r photo) Create(ctx context.Context, userID, title, url string) error {
	const query = `INSERT INTO photos (title, user_id, url) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, title, userID, url)
	if err != nil {
		return err
	}

	return nil
}

func (r photo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	const query = `DELETE FROM photos WHERE id = $1`

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r photo) GetByID(ctx context.Context, id string) (*domain.Photo, error) {
	const query = `SELECT * FROM photos WHERE id = $1 FETCH FIRST 1 ROW ONLY`

	p := &domain.Photo{}
	err := r.db.GetContext(ctx, p, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrPhotoNotFound
		default:
			return nil, err
		}
	}

	return p, nil
}

func (r photo) GetAll(ctx context.Context) ([]*domain.Photo, error) {
	const query = `SELECT * FROM photos`

	pp := make([]*domain.Photo, 0)
	err := r.db.SelectContext(ctx, &pp, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return pp, nil
		default:
			return nil, err
		}
	}
	return pp, nil
}

func (r photo) GetAllByUserID(ctx context.Context, userID string) ([]*domain.Photo, error) {
	const query = `SELECT * FROM photos WHERE user_id = $1`

	pp := make([]*domain.Photo, 0)
	err := r.db.SelectContext(ctx, &pp, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return pp, nil
		default:
			return nil, err
		}
	}
	return pp, nil
}
