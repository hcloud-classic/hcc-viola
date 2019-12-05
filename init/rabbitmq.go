package init

import (
	"hcc/viola/action/rabbitmq"
	"hcc/viola/lib/logger"
	"time"
)

func rabbitmqInit() error {

	for i := 0; i < 100; i++ {
		err := rabbitmq.PrepareChannel()
		if err != nil {
			logger.Logger.Println(err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}

	for i := 0; i < 10; i++ {
		err := rabbitmq.ConsumeAction()
		if err != nil {
			logger.Logger.Println(err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			break
		}
	}


	forever := make(chan bool)
	logger.Logger.Println("RabbitMQ forever channel ready.")
	<-forever

	return nil
}
