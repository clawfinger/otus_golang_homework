package calendarconfig

import (
	"github.com/spf13/viper"
)

type Data struct {
	Logger LoggerConf
	Grpc   Grpc
}

func newConfigData() *Data {
	return &Data{}
}

func (d *Data) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Grpc.SetDefault(v)
}

type LoggerConf struct {
	Level    string
	Filename string
}

func (d *LoggerConf) SetDefault(v *viper.Viper) {
	v.SetDefault("Logger", map[string]interface{}{
		"Level":    "info",
		"Filename": "calendar.log",
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
