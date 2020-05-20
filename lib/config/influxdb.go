package config

type influxdb struct {
	IP   string `goconf:"influxdb:influxdb_ip"`   // ID : InfluxDB  IP
	Port string `goconf:"influxdb:influxdb_port"` // Port : InfluxDB Port
}

// InfluxDB : influxdb config structure
var InfluxDB influxdb
