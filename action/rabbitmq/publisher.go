package rabbitmq

import (
	"encoding/json"
	"hcc/viola/lib/logger"
	"hcc/viola/model"

	"github.com/streadway/amqp"
)

// ViolaToViolin : Provide Some Action to violin
func ViolaToViolin(action model.Control) error {
	qCreate, err := Channel.QueueDeclare(
		"viola_to_violin",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("ViolaToViolin: Failed to declare a create queue")
		return err
	}

	body, _ := json.Marshal(action)
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
		logger.Logger.Println("ViolaToViolin: Failed to register publisher")
		return err
	}

	return nil
}
