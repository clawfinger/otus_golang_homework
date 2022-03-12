package config

import (
	"github.com/spf13/viper"
)

type ConfigData struct {
	DbData  DatabaseData
	Logger  LoggerConf
	Storage Storage
	Http    Http
}

func newConfigData() *ConfigData {
	return &ConfigData{}
}

func (d *ConfigData) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Storage.SetDefault(v)
	d.Http.SetDefault(v)
}

type Http struct {
	Addr string
}

func (d *Http) SetDefault(v *viper.Viper) {
	v.SetDefault("Http", map[string]interface{}{
		"Addr": "127.0.0.1:8080",
	})
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

type DatabaseData struct {
	Username string
	Password string
}

type Storage struct {
	Type string
}

func (d *Storage) SetDefault(v *viper.Viper) {
	v.SetDefault("storage", map[string]interface{}{
		"Type": "inmemory",
	})
}
