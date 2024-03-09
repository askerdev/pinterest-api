package application

import "github.com/gofiber/fiber/v3"

var ServerError = fiber.Error{
	Message: "server error",
	Code:    fiber.StatusInternalServerError,
}
