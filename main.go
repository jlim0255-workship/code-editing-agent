package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	// "github.com/anthropics/anthropic-sdk-go"
	"github.com/ollama/ollama/api"
	"github.com/invopop/jsonschema"
)

type ToolDefinition struct {
	Name string
	Description string
	InputSchema api.ToolFunctionParameters
	Function func(input json.RawMessage) (string, error)
}

type Agent struct {
	client *api.Client
	getUserMessage func() (string, bool)
	tools []ToolDefinition
}

func NewAgent(
	client *api.Client, 
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
	)*Agent {
	return &Agent{
		client: client,
		getUserMessage: getUserMessage,
		tools: tools,
	}
}

// tools
var ReadFileDefinition = ToolDefinition{
	Name: "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function: ReadFile,
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		panic(err)
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GenerateSchema[T any]() api.ToolFunctionParameters {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	// invopop/jsonschema and Ollama use different property-map types. JSON is
	// the common JSON Schema representation, so use it to convert safely.
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		panic(fmt.Sprintf("marshal generated schema: %v", err))
	}

	var parameters api.ToolFunctionParameters
	if err := json.Unmarshal(schemaJSON, &parameters); err != nil {
		panic(fmt.Sprintf("convert generated schema for Ollama: %v", err))
	}

	return parameters
}


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
	tools := []ToolDefinition{ReadFileDefinition}
	agent := NewAgent(client, getUserMessage, tools)
	
	// #fixed err := agent.Run(context.TODO())
	err = agent.Run(context.TODO())

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
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
	// #fixed tools := a.tools
	ollamaTools := make([]api.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		// Convert the generic JSON Schema held by ToolDefinition into the
		// strongly typed schema required by Ollama's Go client.
		schemaJSON, err := json.Marshal(tool.InputSchema)
		if err != nil {
			return nil, fmt.Errorf("marshal schema for tool %q: %w", tool.Name, err)
		}

		var parameters api.ToolFunctionParameters
		if err := json.Unmarshal(schemaJSON, &parameters); err != nil {
			return nil, fmt.Errorf("parse schema for tool %q: %w", tool.Name, err)
		}

		ollamaTools = append(ollamaTools, api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  parameters,
			},
		})
	}

	var response api.Message

	// Define the chat request configuration
	req := &api.ChatRequest{
		Model:    "gpt-oss:120b-cloud",
		Messages: conversation,
		Stream:   new(bool), // Defaults to false, or pass a false pointer explicitly if needed
		Tools:    ollamaTools,
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
