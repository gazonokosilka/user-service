package app

import (
	"context"
	"log/slog"
	"time"
	"user-service/internal/app/rest"
	"user-service/internal/config"
	"user-service/internal/lib/migrator"
	"user-service/internal/storage/psql"
)

type App struct {
	log     *slog.Logger
	storage *psql.Storage
	restApp *rest.App
}

// MustNew constructor of all app components
func MustNew(log *slog.Logger) *App {
	cfg := config.MustLoad()

	storage := psql.Init(cfg.Postgres)

	if err := migrator.RunMigrations(cfg.Postgres, log); err != nil {
		log.Error("failed to run migrations", "error", err)
		panic(err)
	}

	//authRepo := auth.New(storage.GetDB())
	//permRepo := permission.New(storage.GetDB())
	//redisToken := token.New(rdb.Get())
	//codeProvider := code.New(rdb.Get())

	//authService := service.New(
	//	log,
	//	authRepo,
	//	permRepo,
	//	redisToken,
	//	codeProvider,
	//	cfg.TokenTTL,
	//	cfg.SecretKey,
	//)

	restApp := rest.New(
		log,
		//authService,
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

// GracefulShutdown safety stop app
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
