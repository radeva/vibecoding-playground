package main

import (
	"log"
	"net/http"

	"bookstore-api/handlers"
	"bookstore-api/models"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the book store
	store := models.NewBookStore()
	handler := handlers.NewBookHandler(store)

	// Create a new router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/books", handler.CreateBook).Methods("POST")
	router.HandleFunc("/books", handler.GetAllBooks).Methods("GET")
	router.HandleFunc("/books/{id}", handler.GetBook).Methods("GET")
	router.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", handler.DeleteBook).Methods("DELETE")

	// Start the server
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
} 