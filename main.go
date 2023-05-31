package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"textgenie/config"
	"textgenie/rate_limit"
	"time"

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
	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleIncomingSMS(w http.ResponseWriter, r *http.Request, rl *rate_limit.RateLimiter) {
	from := r.FormValue("From")
	body := r.FormValue("Body")

	log.Printf("Received SMS from %s: %s\n", from, body)

	twilioClient := twilio.NewRestClient()

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(config.TwilioPhoneNumber)
	params.SetTo(from)

	if rl.CheckRateLimit(from, 1, time.Minute) == false {
		params.SetBody("Too many requests!")
		twilioClient.Api.CreateMessage(params)

		// We don't return a 429 since we want to let the user know via text
	} else {
		openAiClient := openai.NewClient(config.OpenAiToken)
		resp1, err1 := openAiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: body,
					},
				},
			},
		)

		if err1 != nil {
			fmt.Printf("ChatCompletion error: %v\n", err1)
			return
		}

		fmt.Println(resp1.Choices[0].Message.Content)

		params.SetBody(resp1.Choices[0].Message.Content)

		resp, err := twilioClient.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if resp.Sid != nil {
				fmt.Println(*resp.Sid)
			} else {
				fmt.Println(resp.Sid)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}
