package senderapp

import (
	"context"

	senderconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/sender"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/clawfinger/hw12_13_14_15_calendar/internal/rabbitmq"
)

type Sender struct {
	consumer rabbit.Consumer
	cfg      *senderconfig.Config
	logger   logger.Logger
}

func New(cfg *senderconfig.Config, logger logger.Logger, consumer *rabbit.Consumer) *Sender {
	return &Sender{
		consumer: *consumer,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *Sender) Run(ctx context.Context) {
	go s.consumer.Handle(ctx)
	done := ctx.Done()

	go func() {
		for {
			select {
			case <-done:
				return
			case msg := <-s.consumer.ReadChan:
				s.logger.Info("Got notification from rabbit\n %s", msg)
			}
		}
	}()
	<-ctx.Done()
}
