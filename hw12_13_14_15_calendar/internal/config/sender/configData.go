package senderconfig

import (
	"github.com/spf13/viper"
)

type Data struct {
	Logger   Logger
	Consumer Consumer
}

func newConfigData() *Data {
	return &Data{}
}

func (d *Data) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Consumer.SetDefault(v)
}

type Logger struct {
	Level    string
	Filename string
}

func (d *Logger) SetDefault(v *viper.Viper) {
	v.SetDefault("Logger", map[string]interface{}{
		"Level":    "debug",
		"Filename": "sender.log",
	})
}

type Consumer struct {
	RabbutURL    string
	ExchangeName string
	ExchangeType string
	QueueName    string
}

func (p *Consumer) SetDefault(v *viper.Viper) {
	v.SetDefault("Consumer", map[string]interface{}{
		"RabbutURL":    "127.0.0.1:5672",
		"ExchangeName": "calendarEx",
		"ExchangeType": "topic",
		"QueueName":    "events",
	})
}
