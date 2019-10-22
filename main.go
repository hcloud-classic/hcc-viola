package main

import (
	"fmt"
	"hcc/viola/lib/config"
	"hcc/viola/lib/logger"
	"strings"
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

	// controlcli.HccCli("hcc nodes status 0")
	// // controlcli.NodeInit()
	// controlcli.HccCli("hcc nodes add 2")
	// fmt.Println("Result\n")
	// controlcli.HccCli("hcc nodes status 0")
	// controlcli.HccCli("krgadm nodes add -n 1:2")
	var qwe []string
	qwe = append(qwe, "2:3")
	asd := strings.Split(qwe[0], ":")
	fmt.Println(asd[0], "++++++++++++++++", asd[1])
}
