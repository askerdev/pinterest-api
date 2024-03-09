package main

import (
	_ "github.com/askerdev/pinterest.gallery/docs"
	"github.com/askerdev/pinterest.gallery/internal/application"
	"github.com/askerdev/pinterest.gallery/internal/handler"
	"github.com/askerdev/pinterest.gallery/internal/repository"
	"github.com/askerdev/pinterest.gallery/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

// @title Gallery Microservice
// @version 0.1

// @BasePath /
// @host localhost:8082

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	application.MustSetupConfig()

	fx.New(
		// application related
		fx.Provide(
			application.NewLogger,
			application.NewDatabase,
			application.NewHttpServer,
		),
		// repositories
		fx.Provide(
			fx.Annotate(
				repository.NewLikes,
				fx.As(new(repository.Likes)),
			),
			fx.Annotate(
				repository.NewPhoto,
				fx.As(new(repository.Photo)),
			),
		),
		// services
		fx.Provide(
			service.NewPhoto,
		),
		// handlers
		fx.Invoke(
			handler.RegisterPhoto,
		),
	).Run()
}
