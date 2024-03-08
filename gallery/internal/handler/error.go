package handler

import "github.com/gofiber/fiber/v3"

var ErrInvalidJson = fiber.NewError(fiber.StatusUnprocessableEntity, "invalid json")
