package rabbitmq

import (
	"hcc/viola/lib/logger"
	"log"
)

// XXX : xxx
func XXX() error {
	qCreate, err := Channel.QueueDeclare(
		"get_nodes",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("get_nodes: Failed to declare a create queue")
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
		logger.Logger.Println("get_nodes: Failed to register consumer")
		return err
	}

	go func() {
		for d := range msgsCreate {

			log.Printf("XXX: Received a xxx message: %s", d.Body)

			//logger.Logger.Println("get_nodes: publishing 'create_volume' to cello of server UUID = " + serverUUID)

		}
	}()

	return nil
}
