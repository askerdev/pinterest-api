package domain

import "github.com/gofiber/fiber/v3"

type Photo struct {
	ID     string `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Url    string `json:"url" db:"url"`
	UserID string `json:"user_id" db:"user_id"`
}

var ErrPhotoNotFound = fiber.NewError(fiber.StatusNotFound, "photo not found")
