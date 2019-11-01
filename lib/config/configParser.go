package config

import (
	"hcc/viola/lib/logger"

	"github.com/Terry-Mao/goconf"
)

var conf = goconf.New()
var config = violaConfig{}
var err error

func parseHTTP() {
	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseRabbitMQ() {
	config.RabbitMQConfig = conf.Get("rabbitmq")
	if config.RabbitMQConfig == nil {
		logger.Logger.Panicln("no rabbitmq section")
	}

	RabbitMQ = rabbitmq{}
	RabbitMQ.ID, err = config.RabbitMQConfig.String("rabbitmq_id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Password, err = config.RabbitMQConfig.String("rabbitmq_password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Address, err = config.RabbitMQConfig.String("rabbitmq_address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	RabbitMQ.Port, err = config.RabbitMQConfig.Int("rabbitmq_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseViola() {
	config.RabbitMQConfig = conf.Get("viola")
	if config.RabbitMQConfig == nil {
		logger.Logger.Panicln("no viola section")
	}

	Viola = viola{}
	Viola.NodeAddRetryCount, err = config.ViolaConfig.String("viola_node_add_retry_count")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Viola.NodeAddWaitSec, err = config.ViolaConfig.Int("viola_node_add_wait_sec")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseHTTP()
	parseRabbitMQ()
	parseViola()
}
