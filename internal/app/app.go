package app

import (
	"log/slog"
	"time"

	grpcapp "sso/internal/app/grpc"
	"sso/internal/lib/sl"
	"sso/internal/services/auth"
	"sso/internal/storage/sqlite"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *slog.Logger, port int, tokenTTL time.Duration, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("faild connect to db", sl.Err(err))
		return nil
	}

	authSevice := auth.New(log, storage, tokenTTL)

	grpcApp := grpcapp.New(log, port, authSevice)

	return &App{
		GRPCApp: grpcApp,
	}
}
