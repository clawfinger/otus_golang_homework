package rabbit

import (
	"context"
	"fmt"
	"time"

	//nolint
	"github.com/cenkalti/backoff/v3"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Consumer struct {
	ReadChan    chan string
	consumerTag string
	queue       string
	bindingKey  string
	connector
}

func NewConsumer(uri string, exchangeName string, exchangeType string,
	queueName string, maxReconTime time.Duration, logger logger.Logger) *Consumer {
	return &Consumer{
		ReadChan:    make(chan string),
		consumerTag: "event_tag",

		queue:      queueName,
		bindingKey: queueName,
		connector: connector{
			uri:          uri,
			exchangeName: exchangeName,
			exchangeType: exchangeType,
			done:         make(chan error),
			logger:       logger,
		},
	}
}

func (c *Consumer) Handle(ctx context.Context) {
	msgs, err := c.announceQueue()
	if err != nil {
		c.logger.Error("Announce queue error: ", err)
	}

	ctxDone := ctx.Done()
	for {
		select {
		case <-ctxDone:
			c.conn.Close()
			return
		case msg := <-msgs:
			payload := string(msg.Body)
			c.channel.Ack(msg.DeliveryTag, false)
			c.ReadChan <- payload
		case err := <-c.done:
			if err != nil {
				msgs, err = c.reConnect(ctx)
				if err != nil {
					c.logger.Error("Reconnecting Error: ", err)
					return
				}
				c.logger.Info("Reconnected... possibly")
			}
		}
	}
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

//nolint
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
