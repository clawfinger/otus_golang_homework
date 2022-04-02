package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/app"
	calendarconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	servers "github.com/clawfinger/hw12_13_14_15_calendar/internal/server"
	grpcserver "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/grpc/server"
	internalhttp "github.com/clawfinger/hw12_13_14_15_calendar/internal/server/http"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/clawfinger/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/clawfinger/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	rootCmd := &cobra.Command{
		Use: "calendar",
		Run: func(cmd *cobra.Command, args []string) {
			config := calendarconfig.NewConfig()
			err := config.Init(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on config init, Reason: %s", err.Error())
				return
			}
			logger, err := logger.New(config.Data.Logger.Level, config.Data.Logger.Filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on logger init, Reason: %s", err.Error())
				return
			}
			storageType := config.Data.Storage.Type
			var abstractStorage storage.Storage

			switch storageType {
			case "inmemory":
				abstractStorage = memorystorage.NewMemoryStorage()
			case "sql":
				sqlStorage := sqlstorage.NewSQLStorage(config, logger)
				err := sqlStorage.Connect(ctx)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on sql storage connect, Reason: %s", err.Error())
					return
				}
				abstractStorage = sqlStorage
			}
			serverCtx := servers.NewServerContext(config, abstractStorage, logger)
			httpServer := internalhttp.NewServer(serverCtx)
			grpcServer := grpcserver.NewGrpcServer(serverCtx)
			defer httpServer.Stop(ctx)
			defer grpcServer.Stop()
			defer abstractStorage.Close(ctx)
			app := app.New(config, logger, abstractStorage, httpServer, grpcServer)

			app.Run(ctx)
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	})
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	defaultCfgPath := path.Join(filepath.Dir(executablePath), "config.json")
	flags := rootCmd.Flags()
	flags.StringVarP(&configFile, "config", "c", defaultCfgPath, "Config file path")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command line error")
	}
}
