package rabbitmq

import (
	"encoding/json"
	"hcc/viola/lib/logger"
	"hcc/viola/model"
	"log"
)

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

		}
	}()

	return nil
}
