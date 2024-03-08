package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/askerdev/pinterest.gallery/internal/domain"
	"github.com/jmoiron/sqlx"
)

const likeTable = `
CREATE TABLE IF NOT EXISTS likes (
	photo_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (photo_id, user_id)
);
`

type Likes interface {
	Create(ctx context.Context, userID, photoID string) error
	DeleteByPhotoAndUserID(ctx context.Context, userID, photoID string) error

	GetCountByPhotoID(ctx context.Context, id string) int
	GetByPhotoAndUserID(ctx context.Context, userID, photoID string) (*domain.Like, error)
}

type likes struct {
	db *sqlx.DB
}

func NewLikes(db *sqlx.DB) Likes {
	db.Exec(likeTable)
	return &likes{db}
}

func (r likes) Create(ctx context.Context, userID, photoID string) error {
	const query = `INSERT INTO likes (user_id, photo_id) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, userID, photoID)
	if err != nil {
		return err
	}

	return nil
}

func (r likes) DeleteByPhotoAndUserID(ctx context.Context, userID, photoID string) error {
	_, err := r.GetByPhotoAndUserID(ctx, userID, photoID)
	if err != nil {
		return err
	}

	const query = `DELETE FROM likes WHERE user_id = $1 AND photo_id = $2`

	_, err = r.db.ExecContext(ctx, query, userID, photoID)
	if err != nil {
		return err
	}

	return nil
}

func (r likes) GetCountByPhotoID(ctx context.Context, id string) int {
	const query = `SELECT count(*) FROM likes WHERE photo_id = $1`

	var count int
	r.db.GetContext(ctx, &count, query, id)

	return count
}

func (r likes) GetByPhotoAndUserID(ctx context.Context, userID, photoID string) (*domain.Like, error) {
	const query = `SELECT * FROM likes WHERE user_id = $1 AND photo_id = $2 FETCH FIRST 1 ROW ONLY`

	like := &domain.Like{}
	err := r.db.GetContext(ctx, like, query, userID, photoID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrLikeNotFound
		default:
			return nil, err
		}
	}

	return like, nil
}
