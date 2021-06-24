package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/viola/viola.conf"

type violaConfig struct {
	HTTPConfig     *goconf.Section
	RabbitMQConfig *goconf.Section
	InfluxDBConfig *goconf.Section
	NetworkConfig  *goconf.Section
}

/*-----------------------------------
         Config File Example

##### CONFIG START #####
[http]
port 7000

[rabbitmq]
rabbitmq_id admin
rabbitmq_password qwe1212!Q
rabbitmq_address 192.168.110.10
rabbitmq_port 5672

[influxdb]
influxdb_ip 192.168.110.10

-----------------------------------*/
