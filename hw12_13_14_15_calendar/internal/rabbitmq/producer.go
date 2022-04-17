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

type Producer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	sendChan    chan string
	queueName   string
	maxInterval time.Duration
	connector
}

func NewProducer(uri string, exchangeName string, exchangeType string,
	queueName string, maxReconTime time.Duration, logger logger.Logger) *Producer {
	return &Producer{
		conn:        nil,
		channel:     nil,
		sendChan:    make(chan string, 10),
		queueName:   queueName,
		maxInterval: maxReconTime,
		connector: connector{
			uri:          uri,
			exchangeName: exchangeName,
			exchangeType: exchangeType,
			done:         make(chan error),
			logger:       logger,
		},
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

//nolint
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
