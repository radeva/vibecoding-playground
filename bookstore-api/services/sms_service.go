package services

import (
	"fmt"
	"log"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// SMSService handles sending SMS notifications
type SMSService struct {
	client *twilio.RestClient
	from   string
	to     string
	enabled bool
}

// NewSMSService creates a new instance of SMSService
func NewSMSService() *SMSService {
	// Get Twilio credentials from environment variables
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	toNumber := os.Getenv("TWILIO_TO_NUMBER")

	// Check if all required environment variables are set
	enabled := accountSid != "" && authToken != "" && fromNumber != "" && toNumber != ""
	
	if !enabled {
		log.Println("SMS service disabled: Missing Twilio credentials in environment variables")
		return &SMSService{
			enabled: false,
		}
	}

	// Create Twilio client
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &SMSService{
		client: client,
		from:   fromNumber,
		to:     toNumber,
		enabled: true,
	}
}

// SendBookAddedNotification sends an SMS when a new book is added
func (s *SMSService) SendBookAddedNotification(title, author string) error {
	// If SMS service is disabled, just log and return
	if !s.enabled {
		log.Printf("SMS notification skipped (service disabled): New book '%s' by %s", title, author)
		return nil
	}

	// Create message
	message := fmt.Sprintf("New book added to the library: '%s' by %s", title, author)

	// Send SMS
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(s.to)
	params.SetFrom(s.from)
	params.SetBody(message)

	// Send the message
	_, err := s.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Error sending SMS: %v", err)
		return err
	}

	log.Printf("SMS sent successfully: %s", message)
	return nil
} 