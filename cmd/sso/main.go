package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// Config
	cfg := config.MustLoad()

	// Logger
	logger := setupLogger(cfg.Env)

	// App
	application := app.New(logger, cfg.GRPC.Port, cfg.TokenTTL, cfg.StoragePath)

	// Run gRPC Server
	go application.GRPCApp.MustRun()

	logger.Debug("Server running")

	stop := make(chan os.Signal, 1)
	// якщо прийде SIGINT або SIGTERM, перешли їх у цей канал
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	signal := <-stop

	logger.Info("server stopping", slog.String("signal", signal.String()))

	application.GRPCApp.Stop()
	logger.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}
