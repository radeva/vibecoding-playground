package handlers

import (
	"net/http"
	"time"

	"bookstore-api/models"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	store *models.BookStore
}

func NewBookHandler(store *models.BookStore) *BookHandler {
	return &BookHandler{store: store}
}

// CreateBook handles the creation of a new book
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a simple ID (in a real app, use UUID)
	book.ID = time.Now().Format("20060102150405")
	h.store.Books[book.ID] = book

	c.JSON(http.StatusCreated, book)
}

// GetBook retrieves a book by ID
func (h *BookHandler) GetBook(c *gin.Context) {
	id := c.Param("id")

	book, exists := h.store.Books[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// GetAllBooks retrieves all books
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books := make([]models.Book, 0, len(h.store.Books))
	for _, book := range h.store.Books {
		books = append(books, book)
	}

	c.JSON(http.StatusOK, books)
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	if _, exists := h.store.Books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = id
	h.store.Books[id] = book

	c.JSON(http.StatusOK, book)
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if _, exists := h.store.Books[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	delete(h.store.Books, id)
	c.Status(http.StatusNoContent)
} 