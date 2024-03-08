package handler

import (
	"common/pkg/middleware"
	"github.com/askerdev/pinterest.user/internal/service"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
)

type user struct {
	userSvc *service.User
}

func RegisterUser(app *fiber.App, userSvc *service.User) {
	h := &user{
		userSvc: userSvc,
	}
	h.Register(app)
}

func (h *user) Register(app *fiber.App) {
	auth := app.Group("/auth")

	// accessible

	auth.Post("/signup", h.Signup)
	auth.Post("/signin", h.Signin)

	// refresh token middleware

	auth.Use(middleware.NewMiddleware(viper.GetString("jwt_refresh_secret")))

	// protected

	auth.Get(
		"/access",
		h.AccessToken,
	)
}

func (h *user) Signup(c fiber.Ctx) error {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind().Body(&request); err != nil {
		return ErrInvalidJson
	}

	err := h.userSvc.Signup(c.Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *user) Signin(c fiber.Ctx) error {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind().Body(&request); err != nil {
		return ErrInvalidJson
	}

	token, err := h.userSvc.Signin(c.Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return c.JSON(response{
		RefreshToken: token,
	})
}

func (h *user) AccessToken(c fiber.Ctx) error {
	subject := middleware.MustSubject(c)

	token, err := h.userSvc.AccessToken(c.Context(), subject)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
