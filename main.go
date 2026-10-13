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
	baseURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		fmt.Printf("Invalid Ollama URL: %s\n", err)
		return
	}
	client := api.NewClient(baseURL, http.DefaultClient)

	scanner := bufio.NewScanner(os.Stdin)

	// fixed signature to match NewAgent signature
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}

		return scanner.Text(), true
	}

	// pass the ollama client pointer
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

	// Define the response callback (Ollama streams by default unless Stream is set false/handled in pieces)
	// Even with Stream: false, Ollama executes this callback block once with the complete answer.
	err := a.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		// #fixed response = resp.Message
		response = resp.Message
		return nil
	})

	// #fixed return response, err
	return &response, err
}