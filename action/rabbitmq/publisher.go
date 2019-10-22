package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"hcc/viola/lib/logger"
)

// XXXX : xxx
func XXXX(xxx interface{}) error {
	qCreate, err := Channel.QueueDeclare(
		"return_nodes",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("return_nodes: Failed to declare a create queue")
		return err
	}

	body, _ := json.Marshal(xxx)
	err = Channel.Publish(
		"",
		qCreate.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:     "text/plain",
			ContentEncoding: "utf-8",
			Body:            body,
		})
	if err != nil {
		logger.Logger.Println("return_nodes: Failed to register publisher")
		return err
	}

	return nil
}
