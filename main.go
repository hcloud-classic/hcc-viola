package main

import (
	"hcc/viola/action/rabbitmq"
	"hcc/viola/lib/config"
	"hcc/viola/lib/controlcli"
	"hcc/viola/lib/logger"
	"hcc/viola/lib/syscheck"
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
	status, err := controlcli.TelegrafCheck()
	logger.Logger.Println(status, err)

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

	err = rabbitmq.ConsumeAction()
	if err != nil {
		logger.Logger.Println(err)
	}

	forever := make(chan bool)

	logger.Logger.Println(" [*] Waiting for messages. To exit press Ctrl+C")
	<-forever

}
