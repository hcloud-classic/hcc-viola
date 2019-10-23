package main

import (
	"hcc/viola/action/rabbitmq"
	"hcc/viola/lib/config"
	"hcc/viola/lib/logger"
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

	err = rabbitmq.RunHccCLI()
	if err != nil {
		logger.Logger.Panic(err)
	}

	forever := make(chan bool)

	logger.Logger.Println(" [*] Waiting for messages. To exit press Ctrl+C")
	<-forever

	// controlcli.HccCli("hcc nodes")
	// // controlcli.NodeInit()
	// controlcli.HccCli("hcc nodes add 2")
	// fmt.Println("Result\n")
	// controlcli.HccCli("hcc nodes status 0")
	// controlcli.HccCli("krgadm nodes add -n 1:2")

}
