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

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	readChan     chan string
	consumerTag  string
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
	logger       logger.Logger
}

func (c *Consumer) Handle(ctx context.Context) error {
	var err error
	if err = c.Connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}

	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	ctxDone := ctx.Done()
	for {
		select {
		case <-ctxDone:
			c.conn.Close()
			return nil
		case msg := <-msgs:
			payload := string(msg.Body)
			c.readChan <- payload
		case err := <-c.done:
			if err != nil {
				msgs, err = c.reConnect(ctx)
				if err != nil {
					c.logger.Error("Reconnecting Error: ", err)
					return err
				}
				c.logger.Info("Reconnected... possibly")
			}
		}
	}
}

func (c *Consumer) Connect() error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		c.logger.Error("dial error", err)
		return err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		c.logger.Error("channel error", err)
		return err
	}

	go func() {
		err := <-c.conn.NotifyClose(make(chan *amqp.Error))
		c.logger.Info("closing: %s", err)
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel Closed")
	}()

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		c.logger.Error("exchange declare error", err)
		return err
	}

	return nil
}

func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		if err != nil {
			c.logger.Error("queue Declare error", err)
			return nil, err
		}
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		c.logger.Error("setting qos error", err)
		return nil, err
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		c.logger.Error("queue Bind error", err)
		return nil, err
	}

	msgs, err := c.channel.Consume(
		queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("queue Consume error", err)
		return nil, err
	}

	return msgs, nil
}

func (c *Consumer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
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
			c.logger.Info("stop reconnecting")
			return nil, err
		}

		time.Sleep(d)

		if err := c.Connect(); err != nil {
			c.logger.Info("could not connect in reconnect call", err)
			continue
		}
		msgs, err := c.announceQueue()
		if err != nil {
			c.logger.Info("Couldn't connect", err)
			continue
		}

		return msgs, nil
	}
}
