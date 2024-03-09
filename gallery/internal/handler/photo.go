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

type CreatePhotoRequest struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

// Create this function create new photo
// @Summary create new photo
// @Security ApiKeyAuth
// @Accept  json
// @Produce json
// @Param input body CreatePhotoRequest true "photo info"
// @Success 201 {string} Created
// @Failure 400 {object} errs.FormError
// @Failure 500 {object} GenericError
// @Router /photo [post]
func (h *Photo) Create(c fiber.Ctx) error {
	user := middleware.MustUser(c)

	var request CreatePhotoRequest
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

// Delete this function delete photo
// @Summary delete photo
// @Security ApiKeyAuth
// @Produce json
// @Param   id   path    string true "Photo id"
// @Success 200 {string} OK
// @Failure 404 {object} GenericError
// @Failure 500 {object} GenericError
// @Router /photo/{id} [delete]
func (h *Photo) Delete(c fiber.Ctx) error {
	id := c.Params("id")

	err := h.photoSvc.DeleteById(c.Context(), id)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// PhotoFeed this function get photo feed
// @Summary get photo feed
// @Security ApiKeyAuth
// @Produce json
// @Param user_id query string false "get photo feed by user id"
// @Success 200 {array}  domain.Photo
// @Failure 500 {object} GenericError
// @Router /photo [get]
func (h *Photo) PhotoFeed(c fiber.Ctx) error {
	userID := c.Query("user_id")
	if len(userID) < 10 {
		photos, err := h.photoSvc.GetPhotoFeed(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(photos)
	} else {
		photos, err := h.photoSvc.GetUserPhotoFeed(c.Context(), userID)
		if err != nil {
			return err
		}

		return c.JSON(photos)
	}
}
