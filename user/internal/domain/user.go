package domain

import (
	"github.com/gofiber/fiber/v3"
)

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	PhotoUrl string `json:"photo_url" db:"photo_url"`
}

var ErrInvalidPassword = fiber.NewError(fiber.StatusBadRequest, "invalid password")
var ErrUserAlreadyExists = fiber.NewError(fiber.StatusUnprocessableEntity, "user already exists")
var ErrUserNotFound = fiber.NewError(fiber.StatusNotFound, "user not found")
