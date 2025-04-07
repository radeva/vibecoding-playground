# Bookstore API

A simple REST API for managing books in a bookstore, built with Go and Gin.

## Features

- CRUD operations for books
- Asynchronous SMS notifications using Kafka (segment/kafka-go)
- Environment variable configuration with .env file

## Prerequisites

- Go 1.21 or higher
- Apache Kafka
- Twilio account (for SMS notifications)

## Getting Started

1. Clone the repository

2. Start Kafka:

   ```bash
   # Using Docker (recommended)
   docker run -d --name kafka \
     -p 9092:9092 \
     -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
     -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 \
     -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
     wurstmeister/kafka
   ```

3. Copy the `.env.example` file to `.env` and update with your credentials:

   ```bash
   cp .env.example .env
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

- `POST /books` - Create a new book
- `GET /books` - Get all books
- `GET /books/:id` - Get a specific book
- `PUT /books/:id` - Update a book
- `DELETE /books/:id` - Delete a book

## Message Flow

1. When a new book is added via the API, a message is sent to the Kafka topic "new-books" using segment/kafka-go
2. A Kafka consumer reads the message from the topic
3. The consumer sends an SMS notification using Twilio

## Environment Variables

- `TWILIO_ACCOUNT_SID` - Your Twilio account SID
- `TWILIO_AUTH_TOKEN` - Your Twilio auth token
- `TWILIO_FROM_NUMBER` - Your Twilio phone number
- `TWILIO_TO_NUMBER` - Recipient phone number
- `KAFKA_BROKER` - Kafka broker address (default: localhost:9092)
- `PORT` - Server port (default: 8080)
