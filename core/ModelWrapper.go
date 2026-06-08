package core

import (
	"context"
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
	Name        string
	Temperature float64
}
type Message struct {
	Role    string
	Content string
}
type Chat struct {
	Model   Model
	History []components.ChatMessages
}

func (chat *Chat) CreateUserMessage(msg string) {
	message := components.CreateChatMessagesUser(
		components.ChatUserMessage{
			Role:    components.ChatUserMessageRoleUser,
			Content: components.CreateChatUserMessageContentStr(msg),
		},
	)
	chat.History = append(chat.History, message)
}
func (chat *Chat) CreateChatMessage(msg string) {
	message := components.CreateChatMessagesSystem(
		components.ChatSystemMessage{
			Role:    components.ChatSystemMessageRoleSystem,
			Content: components.CreateChatSystemMessageContentStr(msg),
		},
	)
	chat.History = append(chat.History, message)

}
func (chat *Chat) Send(streamer Streamer) string {
	_ = godotenv.Load()

	key := os.Getenv("KEY")
	if key == "" {
		log.Fatal("KEY di openrouter mancante nel .env")
	}

	ctx := context.Background()

	client := openrouter.New(
		openrouter.WithSecurity(key),
	)

	model := chat.Model

	res, err := client.Chat.Send(ctx, components.ChatRequest{
		Model:       new(model.Name),
		Stream:      new(true),
		MaxTokens:   optionalnullable.From(new(int64(150))),
		Temperature: optionalnullable.From(new(model.Temperature)),
		Messages:    chat.History,
	})
	if err != nil {
		log.Fatal(err)
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
	chat.CreateChatMessage(resp)
	return resp
}

func DefaultModel() Model {
	return Model{Name: "openrouter/free", Temperature: 0.5}
}
func Gpt4o() Model {
	return Model{Name: "", Temperature: 0.5}
}

func CodingLowModel() Model {
	return Model{Name: "openai/gpt-oss-20b:free", Temperature: 0.0}
}

func CodingHighModel() Model {
	return Model{Name: "openai/gpt-oss-120b:free", Temperature: 0.0}
}
