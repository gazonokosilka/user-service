package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"user-service/internal/app"
	cslog "user-service/internal/lib/slog"
)

func main() {
	log := slog.New(cslog.NewCustomHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	application := app.MustNew(log)

	go application.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GracefulShutdown()
	log.Info("Application stopped")
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
		),
	)
}
