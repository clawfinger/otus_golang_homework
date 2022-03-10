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
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	_ = ctx
	defer cancel()
	app := app.New()
	rootCmd := &cobra.Command{
		Use: "calendar",
		Run: func(cmd *cobra.Command, args []string) {
			err := app.Init(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on application init, Reason: %s", err.Error())
				return
			}
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
