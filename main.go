package main

import (
	"hcc/viola/action/rabbitmq"
	"hcc/viola/lib/config"
	"hcc/viola/lib/logger"
	"hcc/viola/lib/syscheck"
	"strconv"
	"time"
)

func main() {
	if !syscheck.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	logger.Logger.Println("Will start viola operations after " + strconv.Itoa(int(config.Viola.StartingDelayTimeSec)) + " seconds")
	time.Sleep(time.Second * time.Duration(config.Viola.StartingDelayTimeSec))

	for i := 0; i < 100; i++ {
		err := rabbitmq.PrepareChannel()
		if err != nil {
			logger.Logger.Println(err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	defer func() {
		_ = rabbitmq.Channel.Close()
	}()
	defer func() {
		_ = rabbitmq.Connection.Close()
	}()

	err := rabbitmq.ViolinToViola()
	if err != nil {
		logger.Logger.Println(err)
	}

	forever := make(chan bool)

	logger.Logger.Println(" [*] Waiting for messages. To exit press Ctrl+C")
	<-forever

	// controlcli.HccCli("hcc nodes status 0")
	// // controlcli.NodeInit()
	// controlcli.HccCli("hcc nodes add 2")
	// fmt.Println("Result\n")
	// controlcli.HccCli("hcc nodes status 0")
	// controlcli.HccCli("krgadm nodes add -n 1:2")
}
