package handler

import (
	"encoding/json"
	"go_msg_broker/models"
	"go_msg_broker/utils"
	"net/http"

	"github.com/segmentio/kafka-go"
)

type PostHandler struct {
	kafkaWriter *kafka.Writer
}

func NewPostHandler(kafkaWriter *kafka.Writer) *PostHandler {
	return &PostHandler{
		kafkaWriter: kafkaWriter,
	}
}

func (p *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ResponsWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// publish post to kafka
	msg, err := json.Marshal(req)
	if err != nil {
		utils.ResponsWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = p.kafkaWriter.WriteMessages(r.Context(), kafka.Message{
		Key:   []byte("create_post"),
		Value: msg,
	})

	if err != nil {
		utils.ResponsWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, req)
}
