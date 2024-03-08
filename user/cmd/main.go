package main

import (
	"github.com/askerdev/pinterest.user/internal/application"
	"github.com/askerdev/pinterest.user/internal/application/config"
	"github.com/askerdev/pinterest.user/internal/handler"
	"github.com/askerdev/pinterest.user/internal/repository"
	"github.com/askerdev/pinterest.user/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

func main() {
	config.MustSetup()

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
				repository.NewUser,
				fx.As(new(repository.User)),
			),
		),
		// services
		fx.Provide(
			fx.Annotate(
				service.NewHasher,
				fx.As(new(service.Hasher)),
			),
			service.NewUser,
		),
		// handlers
		fx.Invoke(
			handler.RegisterUser,
		),
	).Run()
}
