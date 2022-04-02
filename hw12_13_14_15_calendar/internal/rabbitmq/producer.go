package rabbit

import (
	"context"
	"time"

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
}

func (p *Producer) Handle(ctx context.Context) error {
	return nil
}

func (p *Producer) Connect() error {
	return nil
}

func (p *Producer) announceQueue() (<-chan amqp.Delivery, error) {
	return nil, nil
}

func (p *Producer) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	return nil, nil
}
