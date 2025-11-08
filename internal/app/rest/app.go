package rest

import (
	"context"
	"log/slog"
	"net/http"

	v1 "user-service/internal/http/v1"
	customerService "user-service/internal/service/customer"

	"github.com/go-chi/chi/v5"
)

type App struct {
	log             *slog.Logger
	customerService *customerService.Service
	httpServer      *http.Server
}

func New(
	log *slog.Logger,
	customerService *customerService.Service,
	port string,
) *App {
	r := chi.NewRouter()

	// Регистрация маршрутов
	v1.SetupRoutes(r, customerService, log)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &App{
		log:             log,
		customerService: customerService,
		httpServer:      httpServer,
	}
}

func (a *App) Run() error {
	const op = "app.rest.Run"
	a.log.With(slog.String("op", op)).Info("starting REST server", "port", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.rest.Stop"
	a.log.With(slog.String("op", op)).Info("stopping REST server")
	return a.httpServer.Shutdown(ctx)
}
