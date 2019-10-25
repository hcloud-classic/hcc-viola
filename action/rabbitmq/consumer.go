package rabbitmq

import (
	"encoding/json"
	"fmt"
	"hcc/viola/lib/controlcli"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
	"log"
)

// GetClusterIP : Consume 'update_subnet' queues from RabbitMQ channel
func GetClusterIP() error {
	qCreate, err := Channel.QueueDeclare(
		"get_cluster_ip",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("get_cluster_ip: Failed to get cluster_ip")
		return err
	}

	msgsCreate, err := Channel.Consume(
		qCreate.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Println("get_cluster_ip: Failed to register cluster_ip")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("get_cluster_ip: Received a create message: %s", d.Body)

			var subnet model.Subnet
			err = json.Unmarshal(d.Body, &subnet)
			if err != nil {
				logger.Logger.Println("update_subnet: Failed to unmarshal cluster_ip data")
				return
			}

			//TODO: queue get_nodes to flute module

			//logger.Logger.Println("update_subnet: UUID = " + subnet.UUID + ": " + result)
		}
	}()

	return nil
}

//RunHccCLI : Hcc Integration of CLI
func RunHccCLI() error {
	qCreate, err := Channel.QueueDeclare(
		"run_hcc_cli",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("RunHccCLI: Failed to get run_hcc_cli")
		return err
	}

	msgsCreate, err := Channel.Consume(
		qCreate.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Logger.Println("RunHccCLI: Failed to register run_hcc_cli")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("RunHccCLI: Received a create message: %s", d.Body)

			var control model.Control
			err = json.Unmarshal(d.Body, &control)
			if err != nil {
				logger.Logger.Println("RunHccCLI: Failed to unmarshal run_hcc_cli data")
				// return
			}
			fmt.Println("RabbitmQ : ", control)
			status, err := controlcli.HccCli(control.HccCommand, control.HccIPRange)
			if !status && err != nil {
				logger.Logger.Println("RunHccCLI: Faild execution command [", control.HccCommand, "]")
				control.HccCommand = "cluster failed"
			} else {
				logger.Logger.Println("RunHccCLI: Success execution command [", control.HccCommand, "]")
				control.HccCommand = "running"
			}
			ProvideViola(control)

			//TODO: queue get_nodes to flute module

			//logger.Logger.Println("update_subnet: UUID = " + subnet.UUID + ": " + result)
		}
	}()

	return nil
}
