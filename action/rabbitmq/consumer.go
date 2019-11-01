package rabbitmq

import (
	"encoding/json"
	"hcc/viola/lib/config"
	"hcc/viola/lib/controlcli"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
	"log"
	"strconv"
	"time"
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
			log.Printf("get_cluster_ip: Received a create message: %s\n", d.Body)

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

// ViolinToViola : Hcc Integration of CLI
func ViolinToViola() error {
	qCreate, err := Channel.QueueDeclare(
		"violin_to_viola",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("ViolinToViola: Failed to get run_hcc_cli")
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
		logger.Logger.Println("ViolinToViola: Failed to register run_hcc_cli")
		return err
	}

	go func() {
		for d := range msgsCreate {
			logger.Logger.Printf("ViolinToViola: Received a create message: %s\n", d.Body)

			var control model.Control
			err = json.Unmarshal(d.Body, &control)
			if err != nil {
				logger.Logger.Println("ViolinToViola: Failed to unmarshal run_hcc_cli data")
				// return
			}

			var i = 0
			for ; i < int(config.Viola.NodeAddRetryCount); i++ {
				logger.Logger.Println("RabbitmQ : ", control)
				err := controlcli.HccCli(control.HccCommand, control.HccIPRange)
				if err != nil {
					logger.Logger.Println("ViolinToViola: Faild execution command [", control.HccCommand, "]")
					control.HccCommand = "cluster failed"

					logger.Logger.Println("ViolinToViola: Retry after " + strconv.Itoa(int(config.Viola.NodeAddRetryWaitSec)) + " second(s)")
					logger.Logger.Println("ViolinToViola: Retry count (" + strconv.Itoa(i+1) + "/" + strconv.Itoa(int(config.Viola.NodeAddRetryCount)) + ")")
					time.Sleep(time.Second * time.Duration(config.Viola.NodeAddRetryWaitSec))

					continue
				} else {
					logger.Logger.Println("ViolinToViola: Success execution command [", control.HccCommand, "]")
					control.HccCommand = "running"
				}

				err = ViolaToViolin(control)
				if err != nil {
					logger.Logger.Println(err)
				}

				break
			}

			if i > int(config.Viola.NodeAddRetryCount) {
				logger.Logger.Println("ViolinToViola: Retry count exceeded")
			}

			//TODO: queue get_nodes to flute module

			//logger.Logger.Println("update_subnet: UUID = " + subnet.UUID + ": " + result)
		}
	}()

	return nil
}
