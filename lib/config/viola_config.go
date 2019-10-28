package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/viola/viola.conf"

type violaConfig struct {
	HTTPConfig     *goconf.Section
	RabbitMQConfig *goconf.Section
	ViolinConfig   *goconf.Section
}

/*-----------------------------------
         Config File Example

##### CONFIG START #####
[mysql]
id user
password pass
address 111.111.111.111
port 9999
database db_name

[http]
port 8888

[rabbitmq]
rabbitmq_id user
rabbitmq_password pass
rabbitmq_address 555.555.555.555
rabbitmq_port 15672
-----------------------------------*/
