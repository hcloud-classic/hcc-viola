package config

import (
	"hcc/viola/lib/logger"
	"os/exec"
	"strings"

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
	InfluxDB.IP = MasterAddr
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

	RabbitMQ.Address = MasterAddr
	if RabbitMQ.Address != nil {
		logger.Logger.Panicln("Node IP nill")
	}
	RabbitMQ.Port, err = config.RabbitMQConfig.Int("rabbitmq_port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseMasterAddr() {
	config.NetworkConfig = conf.Get("network")
	if config.NetworkConfig == nil {
		logger.Logger.Panicln("no network section")
	}

	NetworkConfig.InterfaceName, err = config.NetworkConfig.String("interface_name")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	cmdString := "cat /var/lib/dhclient/$(ls /var/lib/dhclient/ | grep " + NetworkConfig.InterfaceName + " ) | grep -m 1 'routers '|awk '{print $3}' | tr -d ';'"
	cmd := exec.Command("bash", "-c", cmdString)
	cmdout, err := cmd.CombinedOutput()
	if err != nil {
		logger.Logger.Println("Node status error occurred!!")
	} else {
		MasterAddr = strings.TrimSpace(string(cmdout))
	}
}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}
	parseMasterAddr()
	parseHTTP()
	parseRabbitMQ()
	parseInfluxDB()
}
