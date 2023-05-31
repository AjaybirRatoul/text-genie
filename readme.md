# TextGenie

TextGenie is an application that allows users to send text messages that interact with OpenAI GPT (Generative Pre-trained Transformer) and receive a response back. It integrates with Twilio for sending and receiving SMS messages and utilizes the OpenAI GPT API for generating AI-powered responses.

TextGenie incorporates rate limiting using Redis, which is seamlessly integrated into the application. Redis provides a robust and efficient solution for managing and enforcing rate limits, ensuring the stability and reliability of the system.

## Prerequisites

Before running the application, make sure you have the following:

- Docker: Install Docker on your machine from [https://www.docker.com/](https://www.docker.com/).

## Installation

1. Clone the repository:

```
git clone https://github.com/ajaybirratoul/text-genie.git
```

2. Create a new file named `.env` in the project root directory and add the following environment variables:

```plaintext
TWILIO_ACCOUNT_SID=your_twilio_account_sid
TWILIO_AUTH_TOKEN=your_twilio_auth_token
TWILIO_PHONE_NUMBER=your_twilio_phone_number
OPEN_AI_TOKEN=your_openai_token
```

Replace `your_twilio_account_sid`, `your_twilio_auth_token`, `your_twilio_phone_number`, and `your_openai_token` with your respective values.

## Usage

1. Build and start the Docker containers:

```
docker-compose up
```

2. The application will start and be accessible at http://localhost:8080.

3. Expose the application to the internet (e.g., using ngrok) to receive incoming SMS messages.

4. Send an SMS message to the phone number associated with the application.

5. The application will receive the message and send it to the OpenAI GPT API for processing.

6. OpenAI GPT will generate a response based on the message content.

7. The application will send the generated response back to the sender's phone number.

## License

This project is licensed under the MIT License.
