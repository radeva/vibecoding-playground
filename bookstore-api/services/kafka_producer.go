package services

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"bookstore-api/models"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer() (*KafkaProducer, error) {
	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	new_books_topic := os.Getenv("NEW_BOOKS_TOPIC")
	if broker == "" {
		broker = "localhost:9092" // default broker
	}

	// Create Kafka writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    new_books_topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
		topic:  new_books_topic,
	}, nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

func (kp *KafkaProducer) SendBookMessage(book models.Book) error {
	// Convert book to JSON
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return err
	}

	// Create message
	msg := kafka.Message{
		Key:   []byte(book.ID),
		Value: bookJSON,
	}

	// Send message
	err = kp.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Printf("Book message sent for book ID: %s\n", book.ID)
	return nil
} 