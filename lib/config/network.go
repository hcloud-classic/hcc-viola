package config

type network struct {
	InterfaceName string `goconf:"network:interface_name"` // InterfaceName : interface_name ex) eth0

}

// NetworkConfig : network settings
var NetworkConfig network
