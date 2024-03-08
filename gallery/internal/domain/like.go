package domain

import "github.com/gofiber/fiber/v3"

type Like struct {
	PhotoID string `json:"photo_id" db:"photo_id"`
	UserID  string `json:"user_id" db:"user_id"`
}

var ErrLikeNotFound = fiber.NewError(fiber.StatusNotFound, "like not found")
