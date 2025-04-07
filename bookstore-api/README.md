# Bookstore API

A simple REST API for managing books in a bookstore, built with Go and Gin.

## Features

- CRUD operations for books
- SMS notifications when new books are added
- Environment variable configuration with .env file

## Prerequisites

- Go 1.21 or higher
- Twilio account (for SMS notifications)

## Getting Started

1. Clone the repository
2. Copy the `.env.example` file to `.env` and update with your Twilio credentials:
   ```bash
   cp .env.example .env
   ```
3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

- `POST /books` - Create a new book
- `GET /books` - Get all books
- `GET /books/:id` - Get a specific book
- `PUT /books/:id` - Update a book
- `DELETE /books/:id` - Delete a book

## Environment Variables

- `TWILIO_ACCOUNT_SID` - Your Twilio account SID
- `TWILIO_AUTH_TOKEN` - Your Twilio auth token
- `TWILIO_FROM_NUMBER` - Your Twilio phone number
- `TWILIO_TO_NUMBER` - Recipient phone number
- `PORT` - Server port (default: 8080)
