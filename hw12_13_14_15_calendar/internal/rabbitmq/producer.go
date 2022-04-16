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
	sendChan     chan string
	done         chan error
	uri          string
	exchangeName string
	exchangeType string
	queueName    string
	maxInterval  time.Duration
	logger       logger.Logger
}

func NewProducer(uri string, exchangeName string, exchangeType string,
	queueName string, maxReconTime time.Duration, logger logger.Logger) *Producer {
	return &Producer{
		conn:         nil,
		channel:      nil,
		sendChan:     make(chan string, 10),
		done:         make(chan error),
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queueName,
		maxInterval:  maxReconTime,
		logger:       logger,
	}
}

func (p *Producer) Send(msg string) {
	p.sendChan <- msg
}

func (p *Producer) Handle(ctx context.Context) {
	var err error
	ctxDone := ctx.Done()
	for {
		select {
		case <-ctxDone:
			p.conn.Close()
			return
		case msg := <-p.sendChan:
			err = p.channel.Publish(
				p.exchangeName,
				p.queueName,
				false,
				false,
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "text/plain",
					ContentEncoding: "",
					Body:            []byte(msg),
					DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
					Priority:        0,
				},
			)
			if err != nil {
				p.logger.Error("Exchange Publish error: ", err)
			}
		case err := <-p.done:
			if err != nil {
				err = p.reConnect(ctx)
				if err != nil {
					p.logger.Error("Reconnecting Error: ", err)
					return
				}
				p.logger.Info("Reconnected... possibly")
			}
		}
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

func (p *Producer) reConnect(ctx context.Context) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	p.maxInterval = 15 * time.Second
	be.MaxInterval = p.maxInterval

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			err := fmt.Errorf("stop reconnecting")
			p.logger.Info("stop reconnecting")
			return err
		}

		<-time.After(d)
		if err := p.Connect(); err != nil {
			p.logger.Info("could not connect in reconnect call", err)
			continue
		}
		return nil
	}
}
