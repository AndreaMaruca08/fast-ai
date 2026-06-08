package core

import (
	"context"
	"errors"
	"log"
	"os"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/models/operations"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
	"github.com/OpenRouterTeam/go-sdk/types/stream"
	"github.com/joho/godotenv"
)

type Streamer interface {
	Stream(*operations.SendChatCompletionRequestResponse) string
}

type Model struct {
	Name string
	Type string
}
type Message struct {
	Role    string
	Content string
}
type Chat struct {
	Model   Model
	History []components.ChatMessages
}

const (
	Normal   string = "openrouter/free"
	CodeLow  string = "openai/gpt-oss-20b:free"
	CodeHigh string = "openai/gpt-oss-120b:free"
)

func (chat *Chat) CreateUserMessage(msg string) {
	message := components.CreateChatMessagesUser(
		components.ChatUserMessage{
			Role:    components.ChatUserMessageRoleUser,
			Content: components.CreateChatUserMessageContentStr(msg),
		},
	)
	chat.History = append(chat.History, message)
}
func (chat *Chat) SetSystemPrompt(prompt string) {
	message := components.CreateChatMessagesSystem(
		components.ChatSystemMessage{
			Role:    components.ChatSystemMessageRoleSystem,
			Content: components.CreateChatSystemMessageContentStr(prompt),
		},
	)

	chat.History = append([]components.ChatMessages{message}, chat.History...)
}
func (chat *Chat) CreateAssistantMessage(msg string) {
	message := components.CreateChatMessagesAssistant(
		components.ChatAssistantMessage{
			Role: components.ChatAssistantMessageRoleAssistant,
			Content: optionalnullable.From(
				openrouter.Pointer(
					components.CreateChatAssistantMessageContentStr(msg),
				),
			),
		},
	)

	chat.History = append(chat.History, message)
}

var client *openrouter.OpenRouter
var ctx context.Context
var configs ConfigFile

func (chat *Chat) Send(streamer Streamer) (string, error) {
	if client == nil {
		_ = godotenv.Load()

		key := os.Getenv("KEY")
		if key == "" {
			return "", errors.New("KEY di openrouter mancante nel .env")
		}
		ctx = context.Background()
		client = openrouter.New(
			openrouter.WithSecurity(key),
		)

		configs = GetConfigFileContent()
	}

	model := chat.Model

	var config Config
	switch model.Name {
	case Normal:
		config = configs.Normal
	case CodeLow:
		config = configs.CodeLow
	case CodeHigh:
		config = configs.CodeHigh
	}

	chat.SetSystemPrompt(config.PrePrompt)

	res, err := client.Chat.Send(ctx, components.ChatRequest{
		Model:       openrouter.Pointer(model.Name),
		Stream:      openrouter.Pointer(true),
		MaxTokens:   optionalnullable.From(openrouter.Pointer[int64](int64(config.MaxTokens))),
		Temperature: optionalnullable.From(openrouter.Pointer(config.Temperature)),
		Messages:    chat.History,
	})
	if err != nil {
		return "", errors.New("Errore durante la chiamata a OpenRouter: " + err.Error())
	}
	if res.EventStream == nil {
		log.Fatalf("mi aspettavo uno stream, ricevuto: %v", res.Type)
	}

	defer func(EventStream *stream.EventStream[operations.SendChatCompletionRequestResponseBody]) {
		err := EventStream.Close()
		if err != nil {
			panic(err)
		}
	}(res.EventStream)

	resp := streamer.Stream(res)
	chat.CreateAssistantMessage(resp)
	return resp, nil
}

func DefaultModel() Model {
	return Model{Name: Normal, Type: "General purpose"}
}

func CodingLowModel() Model {
	return Model{Name: CodeLow, Type: "Coding low"}
}

func CodingHighModel() Model {
	return Model{Name: CodeHigh, Type: "Coding high"}
}
