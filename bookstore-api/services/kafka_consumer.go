package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bookstore-api/models"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	topic  string
	sms    *SMSService
}

func NewKafkaConsumer(sms *SMSService) (*KafkaConsumer, error) {
	// Get Kafka broker from environment variable
	broker := os.Getenv("KAFKA_BROKER")
	new_books_topic := os.Getenv("NEW_BOOKS_TOPIC")
	if broker == "" {
		broker = "localhost:9092" // default broker
	}

	// Create Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     new_books_topic,
		GroupID:   "bookstore-sms-group",
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return &KafkaConsumer{
		reader: reader,
		topic:  new_books_topic,
		sms:    sms,
	}, nil
}

func (kc *KafkaConsumer) StartConsuming() error {
	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a signal channel to handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start consuming in a goroutine
	go func() {
		for {
			// Read a message
			msg, err := kc.reader.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					log.Println("Consumer context canceled, stopping...")
					return
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Process the message
			var book models.Book
			if err := json.Unmarshal(msg.Value, &book); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Send SMS notification
			if err := kc.sms.SendBookAddedNotification(book.Title, book.Author); err != nil {
				log.Printf("Error sending SMS: %v", err)
			} else {
				log.Printf("SMS notification sent for book: %s by %s", book.Title, book.Author)
			}
		}
	}()

	// Wait for a signal to stop
	<-signals
	log.Println("Received shutdown signal, stopping consumer...")
	cancel()
	return nil
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
} 