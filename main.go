package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

var (
	debug = false
)

func getCompletionFromMessages(
	ctx context.Context,
	client *openai.Client,
	messages []openai.ChatCompletionMessage,
	model string,
) (openai.ChatCompletionResponse, error) {
	if model == "" {
		model = openai.GPT3Dot5Turbo // see another model options: https://platform.openai.com/docs/models
	}

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
		},
	)
	return resp, err
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Cannot load env file. Err: %s", err)
	}
	ctx := context.Background()
	client := openai.NewClient(os.Getenv("KEY"))
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a friendly chatbot.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "Hi, my name is Budi",
		},
	}
	model := openai.GPT3Dot5Turbo

	fmt.Println("Starting chatbot...")
	fmt.Printf("%s: %s\n", messages[0].Role, messages[0].Content)
	fmt.Printf("%s: %s\n", messages[1].Role, messages[1].Content)
	var userInput string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if userInput != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			})
			userInput = ""
		}
		resp, err := getCompletionFromMessages(ctx, client, messages, model)
		if err != nil {
			fmt.Printf("error when sending request to API: %s\n", err)
			return
		}
		if debug {
			fmt.Printf(
				"ID: %s. Created: %d. Model: %s. Choices: %v.\n",
				resp.ID, resp.Created, resp.Model, resp.Choices,
			)
		}

		answer := openai.ChatCompletionMessage{
			Role:    resp.Choices[0].Message.Role,
			Content: resp.Choices[0].Message.Content,
		}
		messages = append(messages, answer)
		fmt.Printf("%s: %s\n", answer.Role, answer.Content)

		fmt.Printf("%s: ", openai.ChatMessageRoleUser)
		scanner.Scan()
		userInput = scanner.Text()
		if scanner.Err() != nil {
			fmt.Printf("error when scanning user input: %s\n", scanner.Err())
		}
		if len(userInput) < 1 {
			break
		}
	}
}
