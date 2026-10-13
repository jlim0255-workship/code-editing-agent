package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	// "github.com/anthropics/anthropic-sdk-go"
	"github.com/ollama/ollama/api"
)

func main() {
	// client := anthropic.NewClient()
	// #fixed client = api.NewClient(url.Parse("http://localhost:11434"), http.DefaultClient)
	baseURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		fmt.Printf("Invalid Ollama URL: %s\n", err)
		return
	}
	client := api.NewClient(baseURL, http.DefaultClient)

	scanner := bufio.NewScanner(os.Stdin)

	// fixed signature to match NewAgent signature
	// #fixed getUserMessage := func(string, bool){
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}

		return scanner.Text(), true
	}

	// pass the ollama client pointer
	// #fixed agent := NewAgent(&client, getUserMessage)
	agent := NewAgent(client, getUserMessage)
	// #fixed err := agent.Run(context.TODO())
	err = agent.Run(context.TODO())

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func NewAgent(client *api.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

type Agent struct {
	client         *api.Client
	getUserMessage func() (string, bool)
}

func (a *Agent) Run(ctx context.Context) error {
	// conversation := []anthropic.MessageParam{}
	// maintain conversation
	conversation := []api.Message{}

	fmt.Println("Chat with gpt-oss:120b-cloud (use ctrl + c to quit)")

	for {
		fmt.Print("\x1b[34mYou\x1b[94m: ")
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		// append the new User message
		// userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInput))
		userMessage := api.Message{
			Role:    "user",
			Content: userInput,
		}
		conversation = append(conversation, userMessage)

		// run ollama inference
		assistantMessage, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}

		// append ollama's response directly to history (it is already on api.Message)
		// #fixed conversation = append(conversation, assistantMessage.ToParam())
		conversation = append(conversation, *assistantMessage)

		// #fixed for _, content := range assistantMessage.Content { ... }
		fmt.Printf("\x1b[35mAssistant\x1b[95m: %s\n", assistantMessage.Content)
	}

	return nil
}

func (a *Agent) runInference(ctx context.Context, conversation []api.Message) (*api.Message, error) {
	var response api.Message

	// Define the chat request configuration
	req := &api.ChatRequest{
		Model:    "gpt-oss:120b-cloud",
		Messages: conversation,
		Stream:   new(bool), // Defaults to false, or pass a false pointer explicitly if needed
	}

	// 6. Define the response callback (Ollama streams by default unless Stream is set false/handled in pieces)
	// Even with Stream: false, Ollama executes this callback block once with the complete answer.
	err := a.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		// #fixed response = resp.Message
		response = resp.Message
		return nil
	})

	// #fixed return response, err
	return &response, err

	// message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
	// 	Model:     anthropic.ModelClaude3_7SonnetLatest,
	// 	MaxTokens: int64(1024),
	// 	Messages:  conversation,
	// })
	// return message, err
}

// ollama endpoints
// Ollama exposes several key endpoints that we’ll work with:

// /api/generate - Generate text based on a prompt (supports streaming)
// /api/chat - Handle multi-turn conversations with message history
// /api/embeddings - Create vector embeddings for text
// /api/tags - List all pulled models

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// func main(){

// 	// request payload
// 	payload := map[string]interface{}{
// 		"model": "gemma4:e2b",
// 		"prompt": "What is the capital of France?",
// 		"stream": true
// 	}

// 	payloadBytes, _ := json.Marshal(payload)

// 	// make the POST request
// 	response, err := http.Post(
// 		"http://localhost:11434/api/generate",
// 		"application/json",
// 		bytes.NewBuffer(payloadBytes),
// 	)

// 	if err != nil {
// 		fmt.Println("Error making request:", err)
// 		return
// 	}

// 	defer response.Body.Close()

// 	// read and parse the response body
// 	body, _ := io.ReadAll(response.Body)
// 	var result map[string]interface{}
// 	json.Unmarshal(body, &result)

// 	fmt.Println("Response:", result["response"])

// }
