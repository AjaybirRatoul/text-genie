package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Global variables
var (
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
	OpenAiToken       string
)

// LoadEnvironmentVariables loads the environment variables from .env file
func LoadEnvironmentVariables() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Assign environment variables to global variables
	TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	TwilioPhoneNumber = os.Getenv("TWILIO_PHONE_NUMBER")
	OpenAiToken = os.Getenv("OPEN_AI_TOKEN")
}
