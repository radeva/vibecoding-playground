package main

import (
	"log"

	"bookstore-api/handlers"
	"bookstore-api/models"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Start the server
	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
} 