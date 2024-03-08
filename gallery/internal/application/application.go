package application

import (
	"context"
	"errors"
	"github.com/askerdev/pinterest.gallery/pkg/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger() *zap.Logger {
	if viper.GetString("app_mode") == "PRODUCTION" {
		return zap.Must(zap.NewProduction())
	}
	return zap.Must(zap.NewDevelopment())
}

func NewFiberErrorHandler(log *zap.Logger) func(c fiber.Ctx, err error) error {
	return func(c fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		var ee errs.ErrMap
		if errors.As(err, &ee) {
			return c.Status(fiber.StatusBadRequest).JSON(errs.FormError{
				Message: "invalid request data",
				Code:    fiber.StatusBadRequest,
				Errors:  ee,
			})
		}

		if code == fiber.StatusInternalServerError {
			log.Error("Internal error", zap.Error(err))
			return c.Status(code).JSON(fiber.Error{
				Message: "server error",
				Code:    code,
			})
		}

		return c.Status(code).JSON(e)
	}
}

func NewHttpServer(lc fx.Lifecycle, log *zap.Logger) *fiber.App {
	http := fiber.New(fiber.Config{
		ErrorHandler: NewFiberErrorHandler(log),
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("starting server")
			go func() { http.Listen(":" + viper.GetString("app_port")) }()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Warn("shutting down server...")
			return http.ShutdownWithContext(ctx)
		},
	})
	return http
}

func NewDatabase(lc fx.Lifecycle) *sqlx.DB {
	db := sqlx.MustConnect("pgx", viper.GetString("postgres_url"))

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
