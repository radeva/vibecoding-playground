package services

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"bookstore-api/models"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	topic  string
	sms    *SMSService
	// Retry configuration
	maxRetries int
	initialBackoff time.Duration
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

	// Get retry configuration from environment variables or use defaults
	maxRetries := 3
	if retries := os.Getenv("KAFKA_MAX_RETRIES"); retries != "" {
		if r, err := strconv.Atoi(retries); err == nil && r > 0 {
			maxRetries = r
		}
	}

	initialBackoff := 1 * time.Second
	if backoff := os.Getenv("KAFKA_INITIAL_BACKOFF"); backoff != "" {
		if b, err := time.ParseDuration(backoff); err == nil && b > 0 {
			initialBackoff = b
		}
	}

	return &KafkaConsumer{
		reader:         reader,
		topic:          new_books_topic,
		sms:            sms,
		maxRetries:     maxRetries,
		initialBackoff: initialBackoff,
	}, nil
}

// processMessageWithRetry attempts to process a message with exponential backoff retry
func (kc *KafkaConsumer) processMessageWithRetry(ctx context.Context, msg kafka.Message) error {
	var book models.Book
	var err error
	
	// First try to unmarshal the message
	if err = json.Unmarshal(msg.Value, &book); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err // Don't retry unmarshaling errors as they're not transient
	}
	
	// Retry sending SMS with exponential backoff
	backoff := kc.initialBackoff
	for attempt := 0; attempt <= kc.maxRetries; attempt++ {
		// Try to send SMS notification
		err = kc.sms.SendBookAddedNotification(book.Title, book.Author)
		if err == nil {
			// Success - log and return
			log.Printf("SMS notification sent for book: %s by %s (attempt %d)", 
				book.Title, book.Author, attempt+1)
			return nil
		}
		
		// Check if context is canceled
		if ctx.Err() != nil {
			return ctx.Err()
		}
		
		// If we've reached max retries, return the error
		if attempt == kc.maxRetries {
			log.Printf("Failed to send SMS after %d attempts: %v", kc.maxRetries+1, err)
			return err
		}
		
		// Log retry attempt
		log.Printf("Retrying SMS notification for book: %s by %s (attempt %d/%d, backoff: %v)", 
			book.Title, book.Author, attempt+1, kc.maxRetries, backoff)
		
		// Wait with exponential backoff
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			// Double the backoff for next attempt
			backoff *= 2
		}
	}
	
	return err
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

			// Process the message with retry logic
			if err := kc.processMessageWithRetry(ctx, msg); err != nil {
				log.Printf("Failed to process message after all retries: %v", err)
				// Here you could implement a dead letter queue if needed
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