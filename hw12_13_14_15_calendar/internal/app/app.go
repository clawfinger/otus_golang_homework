package app

import (
	"context"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/config"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
)

type App struct { // TODO
	Cfg    *config.Config
	Logger logger.Logger
}

type Storage interface { // TODO
}

func New() *App {
	return &App{
		Cfg:    config.NewConfig(),
		Logger: nil,
	}
}

func (a *App) Init(cfgFilePath string) error {
	err := a.Cfg.Init(cfgFilePath)
	if err != nil {
		return err
	}
	a.Logger = logger.New(a.Cfg)
	return nil
}

func (a *App) Run(ctx context.Context) error {
	_ = ctx
	a.Logger.Info("Password: ", a.Cfg.Data.DbData.Password)
	a.Logger.Info("Username: ", a.Cfg.Data.DbData.Username)

	return nil
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
