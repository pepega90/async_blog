package utils

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("error connecting to rabbitmq: %v", err.Error())
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("error opening channel rabbitmq: %v", err.Error())
		return nil
	}
	_, err = ch.QueueDeclare(
		"create_post",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("error creating queue: %v", err.Error())
		return nil
	}

	return ch
}
