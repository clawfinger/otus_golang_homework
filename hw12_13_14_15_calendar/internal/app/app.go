package app

import (
	"context"

	calendarconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/grpc/server"
	internalhttp "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/http"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	Cfg        *calendarconfig.Config
	Logger     logger.Logger
	storage    storage.Storage
	httpServer *internalhttp.Server
	grpcServer *grpcserver.GrpcServer
}

func New(cfg *calendarconfig.Config, logger logger.Logger, storage storage.Storage,
	httpServer *internalhttp.Server, grpcServer *grpcserver.GrpcServer) *App {
	return &App{
		Cfg:        cfg,
		Logger:     logger,
		storage:    storage,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		err := a.grpcServer.Start()
		if err != nil {
			a.Logger.Info("Failed to start grpc server")
		}
	}()
	go func() {
		err := a.httpServer.Start(ctx)
		if err != nil {
			a.Logger.Info("Failed to start http server")
		}
	}()
	<-ctx.Done()
	return nil
}
