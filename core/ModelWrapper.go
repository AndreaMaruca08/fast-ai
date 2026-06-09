package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
	Name string `json:"name"`
	Type string `json:"type"`
}
type Message struct {
	Role    string
	Content string
}
type Chat struct {
	Model   Model                     `json:"model"`
	History []components.ChatMessages `json:"history"`
}

func (chat *Chat) String() string {
	if chat == nil || len(chat.History) == 0 {
		return ""
	}

	var builder strings.Builder

	for _, message := range chat.History {
		switch {
		case message.ChatUserMessage != nil:
			text := userMessageText(message.ChatUserMessage.Content)
			if text == "" {
				continue
			}
			builder.WriteString(WrapIn("Tu:\n", Blue))
			builder.WriteString(text)
			builder.WriteString("\n\n")

		case message.ChatAssistantMessage != nil:
			content, ok := message.ChatAssistantMessage.Content.GetOrZero()
			if !ok {
				continue
			}

			text := assistantMessageText(content)
			if text == "" {
				continue
			}
			builder.WriteString(WrapIn("AI:\n", Green))
			builder.WriteString(FormatText(text))
			builder.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(builder.String())
}

func userMessageText(content components.ChatUserMessageContent) string {
	if content.Str != nil {
		return *content.Str
	}

	return contentItemsText(content.ArrayOfChatContentItems)
}

func assistantMessageText(content components.ChatAssistantMessageContent) string {
	if content.Str != nil {
		return *content.Str
	}

	if content.ArrayOfChatContentItems != nil {
		return contentItemsText(content.ArrayOfChatContentItems)
	}

	if content.Any != nil {
		return fmt.Sprint(content.Any)
	}

	return ""
}

func contentItemsText(items []components.ChatContentItems) string {
	var builder strings.Builder

	for _, item := range items {
		if item.ChatContentText == nil {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(item.ChatContentText.Text)
	}

	return builder.String()
}

const (
	Normal   string = "openrouter/free"
	CodeLow  string = "openai/gpt-oss-120b:free"
	CodeHigh string = "openai/gpt-4o-mini-2024-07-18"
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

func (chat *Chat) requestMessages(systemPrompt string) []components.ChatMessages {
	systemMessage := components.CreateChatMessagesSystem(
		components.ChatSystemMessage{
			Role:    components.ChatSystemMessageRoleSystem,
			Content: components.CreateChatSystemMessageContentStr(systemPrompt),
		},
	)

	messages := make([]components.ChatMessages, 0, len(chat.History)+1)
	messages = append(messages, systemMessage)
	for _, message := range chat.History {
		if message.ChatSystemMessage != nil {
			continue
		}
		messages = append(messages, message)
	}

	return messages
}

func (chat *Chat) removeSystemMessages() {
	history := chat.History[:0]
	for _, message := range chat.History {
		if message.ChatSystemMessage != nil {
			continue
		}
		history = append(history, message)
	}
	chat.History = history
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

	chat.removeSystemMessages()

	res, err := client.Chat.Send(ctx, components.ChatRequest{
		Model:       openrouter.Pointer(model.Name),
		Stream:      openrouter.Pointer(true),
		MaxTokens:   optionalnullable.From(openrouter.Pointer[int64](int64(config.MaxTokens))),
		Temperature: optionalnullable.From(openrouter.Pointer(config.Temperature)),
		Messages:    chat.requestMessages(config.PrePrompt),
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
