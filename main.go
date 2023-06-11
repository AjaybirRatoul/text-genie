package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"textgenie/config"
	"textgenie/rate_limit"

	"github.com/gorilla/mux"
	"github.com/sashabaranov/go-openai"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func main() {
	// Load environment variables
	config.LoadEnvironmentVariables()

	// Check if required environment variables are set
	if config.TwilioAccountSID == "" || config.TwilioAuthToken == "" || config.TwilioPhoneNumber == "" || config.OpenAiToken == "" {
		log.Fatal("One or more required environment variables are missing")
	}

	// Create a new router
	router := mux.NewRouter()

	// Create a new instance of the rate limiter
	limiter := rate_limit.NewRateLimiter()

	// Register the incoming SMS handler
	router.HandleFunc("/incoming-sms", func(w http.ResponseWriter, r *http.Request) {
		handleIncomingSMS(w, r, limiter)
	})

	// Start the server
	addr := ":8080"
	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func handleIncomingSMS(w http.ResponseWriter, r *http.Request, rl *rate_limit.RateLimiter) {
	from := r.FormValue("From")
	body := r.FormValue("Body")

	log.Printf("Received SMS from %s: %s\n", from, body)

	twilioClient := twilio.NewRestClient()

	messageParams := &twilioApi.CreateMessageParams{}
	messageParams.SetFrom(config.TwilioPhoneNumber)
	messageParams.SetTo(from)

	if !rl.CheckRateLimit(from, 1, time.Minute) {
		// We do not return a 429 since we want to let the user know that they are
		// sending too many requests via text
		messageParams.SetBody("Too many requests!")
		twilioClient.Api.CreateMessage(messageParams)
	} else {
		openAiClient := openai.NewClient(config.OpenAiToken)
		completionReq := openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: body,
				},
			},
		}

		completionResp, err := openAiClient.CreateChatCompletion(context.Background(), completionReq)
		if err != nil {
			log.Printf("Error completing chat: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseContent := completionResp.Choices[0].Message.Content
		fmt.Println(responseContent)
		messageParams.SetBody(responseContent)

		messageResp, err := twilioClient.Api.CreateMessage(messageParams)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Println(messageResp.Sid)
		}
	}

	w.WriteHeader(http.StatusOK)
}
