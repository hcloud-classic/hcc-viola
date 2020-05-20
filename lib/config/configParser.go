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

func parseInfluxDB() {
	config.InfluxDBConfig = conf.Get("influxdb")
	if config.InfluxDBConfig == nil {
		logger.Logger.Panicln("no influxdb section")
	}

	InfluxDB = influxdb{}
	InfluxDB.IP, err = config.InfluxDBConfig.String("influxdb_ip")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	InfluxDB.Port, err = config.InfluxDBConfig.String("influxdb_port")
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

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseHTTP()
	parseRabbitMQ()
	parseInfluxDB()
}
