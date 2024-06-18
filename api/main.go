package main

import (
	"context"
	"go_msg_broker/api/handler"
	"go_msg_broker/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	kafkaWriter := utils.ConnectKafka()
	defer kafkaWriter.Close()

	db := utils.ConnectDB()
	defer db.Close(context.Background())

	userHandler := handler.NewUserHandler(kafkaWriter, db)
	postHandler := handler.NewPostHandler(kafkaWriter)

	r.HandleFunc("/users", userHandler.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/posts", postHandler.CreatePost).Methods(http.MethodPost)

	server := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on port :4000")
	log.Fatal(server.ListenAndServe())
}
