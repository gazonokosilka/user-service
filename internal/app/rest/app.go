package rest

import (
	"context"

	"github.com/go-chi/chi/v5"

	"log/slog"
	"net/http"
)

type App struct {
	log *slog.Logger
	// authService *service.Auth
	httpServer *http.Server
	// secret      string
}

func New(
	log *slog.Logger,
	//authService *service.Auth,
	port string,
	// secret string,
) *App {
	r := chi.NewRouter()

	// Route initialising
	// v1.SetupRoutes(r, authService, log, secret)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	return &App{
		log: log,
		//authService: authService,
		httpServer: httpServer,
	}
}

// Run start http server
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
