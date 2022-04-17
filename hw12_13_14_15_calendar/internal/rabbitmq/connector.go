package rabbit

import (
	"errors"

	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type connector struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	uri          string
	exchangeName string
	exchangeType string
	logger       logger.Logger
	done         chan error
}

func (p *connector) Connect() error {
	var err error

	p.conn, err = amqp.Dial(p.uri)
	if err != nil {
		p.logger.Error("dial error", err)
		return err
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		p.logger.Error("channel error", err)
		return err
	}

	go func() {
		err := <-p.conn.NotifyClose(make(chan *amqp.Error))
		p.logger.Info("closing: %s", err)
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		p.done <- errors.New("channel Closed")
	}()

	if err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		p.logger.Error("exchange declare error", err)
		return err
	}

	return nil
}
