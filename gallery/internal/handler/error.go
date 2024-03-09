package handler

import "github.com/gofiber/fiber/v3"

type GenericError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var ErrInvalidJson = fiber.NewError(fiber.StatusUnprocessableEntity, "invalid json")
