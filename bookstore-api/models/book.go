package models

import "time"

// Book represents a book in the bookstore
type Book struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	PublishingDate time.Time `json:"publishing_date"`
}

// BookStore represents our in-memory database
type BookStore struct {
	Books map[string]Book
}

// NewBookStore creates a new instance of BookStore
func NewBookStore() *BookStore {
	return &BookStore{
		Books: make(map[string]Book),
	}
} 