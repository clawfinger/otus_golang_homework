package app

import (
	"context"
	"time"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/config"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/http"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/clawfinger/hw12_13_14_15_calendar/internal/storage/sql"
)

type App struct { // TODO
	Cfg        *config.Config
	Logger     logger.Logger
	storage    storage.Storage
	httpServer *internalhttp.Server
}

func New() *App {
	return &App{
		Cfg:        config.NewConfig(),
		Logger:     nil,
		storage:    nil,
		httpServer: nil,
	}
}

func (a *App) Init(cfgFilePath string) error {
	err := a.Cfg.Init(cfgFilePath)
	if err != nil {
		return err
	}
	a.Logger = logger.New(a.Cfg)

	// a.storage = memorystorage.NewMemoryStorage()
	sql := sqlstorage.NewSqlStorage(a.Cfg)
	sql.Connect(context.Background())
	serverCtx := internalhttp.NewServerContext(a.Cfg, a.storage, a.Logger)
	a.httpServer = internalhttp.NewServer(serverCtx)
	a.storage = sql
	return nil
}

func (a *App) Run(ctx context.Context) error {
	// a.httpServer.Start(ctx)
	event, err := storage.NewEvent("title", time.Now(), 5*time.Minute, "owner")

	a.storage.Create(event)
	res := a.storage.GetEventsForDay(time.Now().Add(time.Hour))
	_ = res
	return err
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
