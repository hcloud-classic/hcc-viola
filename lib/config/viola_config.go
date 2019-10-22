package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/viola/viola.conf"

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

[flute]
flute_server_address 222.222.222.222
flute_server_port 3333
flute_request_timeout_ms 5000

[violin]
violin_server_address 333.333.333.333
violin_server_port 5555
violin_request_timeout_ms 5000

[dhcpd]
dhcpd_local_dhcpd_service_name isc-dhcpd
dhcpd_local_config_file_location /usr/local/etc/dhcpd.conf
dhcpd_config_file_location /etc/harp/dhcpd
dhcpd_min_lease_time 1200
dhcpd_default_lease_time 1800
dhcpd_max_lease_time 3600
-----------------------------------*/
