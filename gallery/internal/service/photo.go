package service

import (
	"context"
	"github.com/askerdev/pinterest.gallery/internal/domain"
	"github.com/askerdev/pinterest.gallery/internal/repository"
)

type Photo struct {
	photoRepo repository.Photo
}

func NewPhoto(photoRepo repository.Photo) *Photo {
	return &Photo{
		photoRepo: photoRepo,
	}
}

func (s *Photo) Create(ctx context.Context, userID, title, url string) error {
	err := s.photoRepo.Create(ctx, userID, title, url)
	if err != nil {
		return err
	}

	return nil
}

func (s *Photo) DeleteById(ctx context.Context, id string) error {
	err := s.photoRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Photo) GetPhotoFeed(ctx context.Context) ([]*domain.Photo, error) {
	photos, err := s.photoRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func (s *Photo) GetUserPhotoFeed(ctx context.Context, userID string) ([]*domain.Photo, error) {
	photos, err := s.photoRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
