package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/viola/viola.conf"

type violaConfig struct {
	HTTPConfig     *goconf.Section
	RabbitMQConfig *goconf.Section
	InfluxDBConfig *goconf.Section
	NetworkConfig  *goconf.Section
}
