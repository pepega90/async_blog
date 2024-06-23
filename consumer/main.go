package main

import (
	"context"
	"encoding/json"
	"go_msg_broker/models"
	"go_msg_broker/utils"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

func main() {
	// kafka consumer
	kafkaReader := utils.ConnectKafkaReader("create_post")
	defer kafkaReader.Close()

	// rabbitmq consumer
	// rabbitCh := utils.ConnectRabbitMQ()
	// defer rabbitCh.Close()

	// nats consumer
	natsConn := utils.ConnectNATS()
	defer natsConn.Close()

	db := utils.ConnectDB()
	defer db.Close(context.Background())

	// consumeKafka(kafkaReader, db)
	// consumeRabbitMQ(rabbitCh, db)
	consumeNats(natsConn, db)
}

func consumeKafka(kafkaReader *kafka.Reader, db *pgx.Conn) {
	log.Println("Kafka consumer is running")

	for {
		msg, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		var post models.Post
		if err := json.Unmarshal(msg.Value, &post); err != nil {
			log.Printf("Error decoding JSON: %s", err)
			continue
		}

		_, err = db.Exec(context.Background(), "INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)", post.UserId, post.Title, post.Content)
		if err != nil {
			log.Fatalf("error insert post to database from kafka: %v", err)
		} else {
			log.Printf("success insert post to database from kafka: %v", post.Title)
		}
	}
}

func consumeRabbitMQ(rabbitCh *amqp091.Channel, db *pgx.Conn) {
	messages, _ := rabbitCh.Consume(
		"create_post",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	go func() {
		for msg := range messages {
			var post models.Post
			if err := json.Unmarshal(msg.Body, &post); err != nil {
				log.Fatalf("error unmarshal json: %v", err.Error())
				break
			}
			_, err := db.Exec(context.Background(), "INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)", post.UserId, post.Title, post.Content)
			if err != nil {
				log.Fatalf("error insert post to database from rabbitmq: %v", err)
			} else {
				log.Printf("success insert post to database from rabbitmq: %v", post.Title)
			}
		}
	}()
	log.Println("[x] Waiting messages, rabbitmq")
	select {}
}

func consumeNats(nc *nats.Conn, db *pgx.Conn) {
	_, err := nc.Subscribe("create_post", func(msg *nats.Msg) {
		var post models.Post
		if err := json.Unmarshal(msg.Data, &post); err != nil {
			log.Fatalf("error unmarshal json: %v", err.Error())
			return

		}

		log.Printf("Received message: %v", post)

		_, err := db.Exec(context.Background(), "INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3)", post.UserId, post.Title, post.Content)
		if err != nil {
			log.Fatalf("error insert post to database from NATS: %v", err)
		} else {
			log.Printf("success insert post to database from NATS: %v", post.Title)
		}
	})
	if err != nil {
		log.Fatalf("error subscribing to nats event: %v", err.Error())
		return
	}
	log.Println("[x] waiting for message from nats")
	select {}
}
