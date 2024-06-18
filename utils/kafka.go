package utils

import "github.com/segmentio/kafka-go"

func ConnectKafka() *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "create_post",
		Balancer: &kafka.LeastBytes{},
	}
}

func ConnectKafkaReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
}
