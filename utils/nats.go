package utils

import (
	"log"

	"github.com/nats-io/nats.go"
)

func ConnectNATS() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("error connecting to nats cluster: %v", err.Error())
		return nil
	}
	return nc
}
