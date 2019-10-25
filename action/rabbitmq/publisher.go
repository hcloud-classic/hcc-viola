package rabbitmq

import (
	"encoding/json"
	"hcc/viola/lib/logger"
	"hcc/viola/model"

	"github.com/streadway/amqp"
)

// ProvideViola : Provide Some Action to violin
func ProvideViola(action model.Control) error {
	qCreate, err := Channel.QueueDeclare(
		"consume_viola",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		logger.Logger.Println("consume_viola: Failed to declare a create queue")
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
		logger.Logger.Println("consume_viola: Failed to register publisher")
		return err
	}

	return nil
}
