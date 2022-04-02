package rabbit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Producer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
	maxInterval  time.Duration
	logger       logger.Logger
}

func (p *Producer) Handle(ctx context.Context) error {
	var err error
	if err = p.Connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	msgs, err := p.announceQueue()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for {
		//do stuff
		_ = msgs
		if <-p.done != nil {
			msgs, err = p.reConnect(ctx)
			if err != nil {
				return fmt.Errorf("reconnecting Error: %s", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

func (p *Producer) Connect() error {
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

func (p *Producer) announceQueue() (<-chan amqp.Delivery, error) {
	queue, err := p.channel.QueueDeclare(
		p.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		if err != nil {
			p.logger.Error("queue Declare error", err)
			return nil, err
		}
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = p.channel.Qos(50, 0, false)
	if err != nil {
		p.logger.Error("setting qos error", err)
		return nil, err
	}

	// Создаём биндинг (правило маршрутизации).
	if err = p.channel.QueueBind(
		queue.Name,
		p.bindingKey,
		p.exchangeName,
		false,
		nil,
	); err != nil {
		p.logger.Error("queue Bind error", err)
		return nil, err
	}

	msgs, err := p.channel.Consume(
		queue.Name,
		p.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		p.logger.Error("queue Consume error", err)
		return nil, err
	}

	return msgs, nil
}

func (p *Producer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			err := fmt.Errorf("stop reconnecting")
			p.logger.Info("stop reconnecting")
			return nil, err
		}

		select {
		case <-time.After(d):
			if err := p.Connect(); err != nil {
				p.logger.Info("could not connect in reconnect call", err)
				continue
			}
			msgs, err := p.announceQueue()
			if err != nil {
				p.logger.Info("Couldn't connect", err)
				continue
			}

			return msgs, nil
		}
	}
}
