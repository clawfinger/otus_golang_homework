package app

import (
	"context"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/config"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/http"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
	Cfg        *config.Config
	Logger     logger.Logger
	storage    storage.Storage
	httpServer *internalhttp.Server
}

func New(cfg *config.Config, logger logger.Logger, storage storage.Storage, httpServer *internalhttp.Server) *App {
	return &App{
		Cfg:        cfg,
		Logger:     logger,
		storage:    storage,
		httpServer: httpServer,
	}
}

func (a *App) Run(ctx context.Context) error {
	err := a.httpServer.Start(ctx)
	return err
}
