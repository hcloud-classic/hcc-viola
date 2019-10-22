package main

import (
	"fmt"
	"hcc/viola/action/rabbitmq"
	"hcc/viola/lib/config"
	"hcc/viola/lib/logger"
	"strings"
)

func main() {
	//if !syscheck.CheckRoot() {
	//	return
	//}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := rabbitmq.PrepareChannel()
	if err != nil {
		logger.Logger.Panic(err)
	}
	defer func() {
		_ = rabbitmq.Channel.Close()
	}()
	defer func() {
		_ = rabbitmq.Connection.Close()
	}()

	// controlcli.HccCLI("nodes", "status")
	wwqwe := strings.Split("krgadm nodes add -n 2", " ")

	fmt.Println(wwqwe[1])

}
