package main

import (
	"context"
	"encoding/json"
	"go_msg_broker/models"
	"go_msg_broker/utils"
	"log"
)

func main() {
	kafkaReader := utils.ConnectKafkaReader("create_post")
	defer kafkaReader.Close()

	db := utils.ConnectDB()
	defer db.Close(context.Background())

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
			log.Fatalf("error insert post to database: %v", err)
		} else {
			log.Printf("success insert post to database: %v", post.Title)
		}
	}
}
