package end

import "hcc/viola/action/rabbitmq"

func rabbitmqEnd() {
	if rabbitmq.Channel != nil {
		_ = rabbitmq.Channel.Close()
	}

	if rabbitmq.Connection != nil {
		_ = rabbitmq.Connection.Close()
	}
}
