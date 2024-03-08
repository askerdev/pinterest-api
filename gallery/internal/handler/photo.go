package handler

import (
	"common/pkg/middleware"
	"github.com/askerdev/pinterest.gallery/internal/service"
	"github.com/askerdev/pinterest.gallery/pkg/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/viper"
	"net/url"
)

type Photo struct {
	photoSvc *service.Photo
}

func RegisterPhoto(app *fiber.App, photoSvc *service.Photo) {
	h := &Photo{photoSvc: photoSvc}

	p := app.Group("/photo")
	p.Use(middleware.NewMiddleware(viper.GetString("jwt_access_secret")))
	p.Post("/", h.Create)
	p.Delete("/:id", h.Delete)
	p.Get("/", h.PhotoFeed)
}

func (h *Photo) Create(c fiber.Ctx) error {
	user := middleware.MustUser(c)

	var request struct {
		Url   string `json:"url"`
		Title string `json:"title"`
	}
	{
		if err := c.Bind().Body(&request); err != nil {
			return ErrInvalidJson
		}

		var errors errs.ErrMap

		if _, err := url.ParseRequestURI(request.Url); err != nil {
			errors.Set("url", "invalid url")
		}
		if len(request.Title) < 3 {
			errors.Set("title", "title too small")
		}
		if len(request.Title) > 60 {
			errors.Set("title", "title too big")
		}
		if errors != nil {
			return errors
		}
	}

	err := h.photoSvc.Create(c.Context(), user.Subject, request.Title, request.Url)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Photo) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	err := h.photoSvc.DeleteById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Photo) PhotoFeed(c fiber.Ctx) error {
	photos, err := h.photoSvc.GetPhotoFeed(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(photos)
}

func (h *Photo) UserPhotoFeed(c fiber.Ctx) error {
	userID := c.Query("user_id")
	photos, err := h.photoSvc.GetUserPhotoFeed(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(photos)
}
