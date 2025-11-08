package app

import (
	"context"
	"log/slog"
	"time"
	"user-service/internal/app/rest"
	"user-service/internal/config"
	"user-service/internal/lib/migrator"
	customerService "user-service/internal/service/customer"
	"user-service/internal/storage/psql"
	customerRepo "user-service/internal/storage/repository/customer"
)

type App struct {
	log     *slog.Logger
	storage *psql.Storage
	restApp *rest.App
}

func MustNew(log *slog.Logger) *App {
	cfg := config.MustLoad()

	storage := psql.Init(cfg.Postgres)

	if err := migrator.RunMigrations(cfg.Postgres, log); err != nil {
		log.Error("failed to run migrations", "error", err)
		panic(err)
	}

	// Инициализация репозитория
	custRepo := customerRepo.New(storage.GetDB())

	// Инициализация сервиса
	custService := customerService.New(log, custRepo)

	restApp := rest.New(
		log,
		custService,
		cfg.Server.Port,
	)

	return &App{
		log:     log,
		storage: storage,
		restApp: restApp,
	}
}

func (a *App) MustRun() {
	const op = "app.MustRun"
	a.log.With(slog.String("op", op)).Info("starting application")

	if err := a.restApp.Run(); err != nil {
		panic(err)
	}
}

func (a *App) GracefulShutdown() {
	const op = "app.GracefulShutdown"
	a.log.With(slog.String("op", op)).Info("shutting down application")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.restApp.Stop(ctx); err != nil {
		a.log.Error("failed to stop HTTP server", err)
	}

	if a.storage != nil {
		a.storage.Close()
		a.log.Info("database connection closed")
	}
}
