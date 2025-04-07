package main

import (
	"fmt"
	"log"
	"os"

	"bookstore-api/handlers"
	"bookstore-api/models"
	"bookstore-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize the book store
	store := models.NewBookStore()

	// Initialize Kafka producer
	producer, err := services.NewKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize SMS service
	smsService := services.NewSMSService()

	// Initialize and start Kafka consumer
	consumer, err := services.NewKafkaConsumer(smsService)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages in a separate goroutine
	go func() {
		if err := consumer.StartConsuming(); err != nil {
			log.Printf("Error consuming messages: %v", err)
		}
	}()

	// Initialize handler with both store and producer
	handler := handlers.NewBookHandler(store, producer)

	// Create a new Gin router
	router := gin.Default()

	// Define routes
	router.POST("/books", handler.CreateBook)
	router.GET("/books", handler.GetAllBooks)
	router.GET("/books/:id", handler.GetBook)
	router.PUT("/books/:id", handler.UpdateBook)
	router.DELETE("/books/:id", handler.DeleteBook)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatal(err)
	}
} 