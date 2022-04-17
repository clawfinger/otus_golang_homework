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

	pb "github.com/clawfinger/hw12_13_14_15_calendar/api/generated"
	schedulerapp "github.com/clawfinger/hw12_13_14_15_calendar/internal/appdata/scheduler"
	schedulerconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/scheduler"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/clawfinger/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configFile string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	rootCmd := &cobra.Command{
		Use: "scheduler",
		Run: func(cmd *cobra.Command, args []string) {
			config := schedulerconfig.NewConfig()
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

			conn, err := grpc.Dial(config.Data.Grpc.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logger.Error("Failed to connect to grpc server")
				return
			}
			grpcsClient := pb.NewCalendarClient(conn)

			rabbitProducer := rabbit.NewProducer(config.Data.Producer.RabbutURL, config.Data.Producer.ExchangeName,
				config.Data.Producer.ExchangeType, config.Data.Producer.QueueName, time.Second*5, logger)

			err = rabbitProducer.Connect()
			if err != nil {
				logger.Error("Error connecting to rabbit", err.Error())
				return
			}
			app := schedulerapp.New(config, logger, grpcsClient, rabbitProducer, config.Data.Scheduler.CycleTime)
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
