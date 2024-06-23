package handler

import (
	"encoding/json"
	"go_msg_broker/models"
	"go_msg_broker/utils"
	"net/http"

	"github.com/nats-io/nats.go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

type PostHandler struct {
	kafkaWriter *kafka.Writer
	rabbitMqPub *amqp.Channel
	natsPub     *nats.Conn
}

func NewPostHandler(kafkaWriter *kafka.Writer, rabbitMqPub *amqp.Channel, natsConnect *nats.Conn) *PostHandler {
	return &PostHandler{
		kafkaWriter: kafkaWriter,
		rabbitMqPub: rabbitMqPub,
		natsPub:     natsConnect,
	}
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponsWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	msg, err := json.Marshal(req)
	if err != nil {
		utils.ResponsWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// publish post to kafka
	// err = p.kafkaWriter.WriteMessages(r.Context(), kafka.Message{
	// 	Key:   []byte("create_post"),
	// 	Value: msg,
	// })

	// publish post to rabbitmq
	// err = p.rabbitMqPub.Publish(
	// 	"",
	// 	"create_post",
	// 	false,
	// 	false,
	// 	amqp.Publishing{
	// 		ContentType: "application/json",
	// 		Body:        msg,
	// 	},
	// )
	// if err != nil {
	// 	utils.ResponsWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	// publish post to NATS
	err = p.natsPub.Publish("create_post", msg)
	if err != nil {
		utils.ResponsWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, req)
}
