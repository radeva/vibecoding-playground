package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bookstore-api/models"

	"github.com/gorilla/mux"
)

type BookHandler struct {
	store *models.BookStore
}

func NewBookHandler(store *models.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// CreateBook handles the creation of a new book
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a simple ID (in a real app, use UUID)
	book.ID = time.Now().Format("20060102150405")
	h.store.Books[book.ID] = book

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// GetBook retrieves a book by ID
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	book, exists := h.store.Books[id]
	if !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// GetAllBooks retrieves all books
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books := make([]models.Book, 0, len(h.store.Books))
	for _, book := range h.store.Books {
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, exists := h.store.Books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = id
	h.store.Books[id] = book

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, exists := h.store.Books[id]; !exists {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	delete(h.store.Books, id)
	w.WriteHeader(http.StatusNoContent)
} 