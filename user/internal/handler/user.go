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

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signup create new account
// @Summary Register an account
// @Tags    auth
// @Accept  json
// @Produce json
// @Param   input  body     SignupRequest true "account info"
// @Success 201    {string} Created
// @Failure 422    {object} GenericError
// @Failure 422    {object} GenericError
// @Failure 500    {object} GenericError
// @Router /auth/signup [post]
func (h *user) Signup(c fiber.Ctx) error {
	var request SignupRequest
	if err := c.Bind().Body(&request); err != nil {
		return ErrInvalidJson
	}

	err := h.userSvc.Signup(c.Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type SigninResponse struct {
	RefreshToken string `json:"refresh_token"`
}

// Signin login into account
// @Summary Login into created account
// @Tags auth
// @Accept  json
// @Produce json
// @Param 	input body     SigninRequest true "account info"
// @Success 200   {object} SigninResponse
// @Failure 422   {object} GenericError
// @Failure 400   {object} GenericError
// @Failure 404   {object} GenericError
// @Failure 500   {object} GenericError
// @Router /auth/signin [post]
func (h *user) Signin(c fiber.Ctx) error {
	var request SigninRequest

	if err := c.Bind().Body(&request); err != nil {
		return ErrInvalidJson
	}

	token, err := h.userSvc.Signin(c.Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return c.JSON(SigninResponse{
		RefreshToken: token,
	})
}

type AccessResponse struct {
	AccessToken string `json:"access_token"`
}

// AccessToken get access token
// @Summary Get new access token (need refresh token)
// @Security ApiKeyAuth
// @Tags auth
// @Produce json
// @Success 200 {object} AccessResponse
// @Failure 401 {object} GenericError
// @Failure 500 {object} GenericError
// @Router /auth/access [get]
func (h *user) AccessToken(c fiber.Ctx) error {
	subject := middleware.MustSubject(c)

	token, err := h.userSvc.AccessToken(c.Context(), subject)
	if err != nil {
		return err
	}

	return c.JSON(AccessResponse{
		AccessToken: token,
	})
}
