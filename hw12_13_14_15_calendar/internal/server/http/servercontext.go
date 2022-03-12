package internalhttp

import (
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/config"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
)

type serverContext struct {
	Cfg     *config.Config
	Storage storage.Storage
	Logger  logger.Logger
}

func NewServerContext(cfg *config.Config, storage storage.Storage, logger logger.Logger) *serverContext {
	return &serverContext{
		Cfg:     cfg,
		Storage: storage,
		Logger:  logger,
	}
}
