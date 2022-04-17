package schedulerconfig

import (
	"time"

	"github.com/spf13/viper"
)

type Data struct {
	Logger    LoggerConf
	Grpc      Grpc
	Producer  Producer
	Scheduler Scheduler
}

func newConfigData() *Data {
	return &Data{}
}

func (d *Data) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Grpc.SetDefault(v)
	d.Producer.SetDefault(v)
	d.Scheduler.SetDefault(v)
}

type LoggerConf struct {
	Level    string
	Filename string
}

func (d *LoggerConf) SetDefault(v *viper.Viper) {
	v.SetDefault("Logger", map[string]interface{}{
		"Level":    "debug",
		"Filename": "scheduler.log",
	})
}

type Grpc struct {
	Addr string
}

func (d *Grpc) SetDefault(v *viper.Viper) {
	v.SetDefault("Grpc", map[string]interface{}{
		"Addr": "127.0.0.1:50051",
	})
}

type Producer struct {
	RabbutURL    string
	ExchangeName string
	ExchangeType string
	QueueName    string
}

func (p *Producer) SetDefault(v *viper.Viper) {
	v.SetDefault("Producer", map[string]interface{}{
		"RabbutURL":    "127.0.0.1:5672",
		"ExchangeName": "calendarEx",
		"ExchangeType": "topic",
		"QueueName":    "events",
	})
}

type Scheduler struct {
	CycleTime time.Duration
}

func (p *Scheduler) SetDefault(v *viper.Viper) {
	v.SetDefault("Scheduler", map[string]interface{}{
		"CycleTime": "1m",
	})
}
