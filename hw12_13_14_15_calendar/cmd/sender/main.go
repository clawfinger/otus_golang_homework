package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	senderapp "github.com/clawfinger/hw12_13_14_15_calendar/internal/appdata/sender"
	senderconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/sender"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/clawfinger/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	rootCmd := &cobra.Command{
		Use: "sender",
		Run: func(cmd *cobra.Command, args []string) {
			config := senderconfig.NewConfig()
			err := config.Init(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on config init, Reason: %s", err.Error())
				return
			}
			logger, err := logger.New(config.Data.Logger.Level, config.Data.Logger.Filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on logger init, Reason: %s", err.Error())
			}
			consumer := rabbit.NewConsumer(config.Data.Consumer.RabbutURL,
				config.Data.Consumer.ExchangeName, config.Data.Consumer.ExchangeType,
				config.Data.Consumer.QueueName, time.Second*5, logger)
			err = consumer.Connect()
			if err != nil {
				logger.Error("Error connecting to rabbit", err.Error())
				return
			}
			app := senderapp.New(config, logger, consumer)
			app.Run(ctx)
		},
	}
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
