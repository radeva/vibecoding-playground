package main

import (
	"fmt"
	"log"
	"os"

	"bookstore-api/handlers"
	"bookstore-api/models"

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
	handler := handlers.NewBookHandler(store)

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