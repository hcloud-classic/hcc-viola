package rabbitmq

import (
	"bytes"
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

//ConsumeAction : Hcc Integration of CLI
func ConsumeAction() error {
	qCreate, err := Channel.QueueDeclare(
		"to_viola",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("ConsumeAction: Failed to get run_hcc_cli")
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
		logger.Logger.Println("ConsumeAction: Failed to register run_hcc_cli")
		return err
	}

	go func() {
		for d := range msgsCreate {
			log.Printf("ConsumeAction: Received a create message: %s", d.Body)

			var control model.Control
			err = json.Unmarshal(d.Body, &control)
			if err != nil {
				logger.Logger.Println("ConsumeAction: Failed to unmarshal run_hcc_cli data")
				// return
			}
			var pretty bytes.Buffer
			error := json.Indent(&pretty, d.Body, "", "\t")
			if error != nil {
				log.Println("JSON parse error: ", error)
				return
			}
			fmt.Println("RabbitmQ : ", control)
			logger.Logger.Println("RabbitmQ : ", string(pretty.Bytes()))
			logger.Logger.Println("Codex : ", control.Control.HccType.HccIPRange)
			status, err := controlcli.HccCli(control)
			errstr := fmt.Sprintf("%v", err)
			if !status && err != nil {
				logger.Logger.Println("ConsumeAction: Faild execution command [", errstr, "]")
				control.Control.ActionResult = "Failed"
			} else {
				logger.Logger.Println("ConsumeAction: Success execution command [", errstr, "]")
				control.Control.ActionResult = "Running"
			}
			logger.Logger.Println("Will Publish Strcut : ", control, "\n To : [", control.Receiver, "]")
			switch control.Receiver {
			case "violin":
				PublishViolin(control)
			// To-Do
			//  if another modules want to receive action result, implementation code write here
			default:
				logger.Logger.Println("No Receiver")
			}

		}
	}()

	return nil
}
