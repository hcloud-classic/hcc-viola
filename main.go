package main

import (
	"hcc/viola/lib/config"
	"hcc/viola/lib/controlcli"
	"hcc/viola/lib/logger"
)

func main() {
	// if !syscheck.CheckRoot() {
	// 	return
	// }

	if !logger.Prepare() {
		return
	}
	defer logger.FpLog.Close()

	config.Parser()

	// err := rabbitmq.PrepareChannel()
	// if err != nil {
	// 	logger.Logger.Panic(err)
	// }
	// defer func() {
	// 	_ = rabbitmq.Channel.Close()
	// }()
	// defer func() {
	// 	_ = rabbitmq.Connection.Close()
	// }()

	// http.Handle("/graphql", graphql.GraphqlHandler)

	// logger.Logger.Println("Server is running on port " + config.HTTPPort)
	// err = http.ListenAndServe(":"+config.HTTPPort, nil)
	// if err != nil {
	// 	logger.Logger.Println("Failed to prepare http server!")
	// }

	// forever := make(chan bool)

	// logger.Logger.Println(" [*] Waiting for messages. To exit press Ctrl+C")
	// <-forever

	controlcli.HccCli("hcc nodes status 0")
	controlcli.nodeOnlineCheck()
}
