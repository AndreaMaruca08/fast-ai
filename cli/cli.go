package cli

import (
	"fast_ai_client/core"
	"fmt"
	"log"
	"strings"

	"github.com/OpenRouterTeam/go-sdk/models/operations"
	"github.com/eiannone/keyboard"
)

func ClearTerminal() {
	fmt.Print("\033[H\033[2J\n")
}

var currentPage *Page

type Streamer struct{}

func (cs *Streamer) Stream(res *operations.SendChatCompletionRequestResponse) string {
	var full strings.Builder
	var formatter Formatter

	for res.EventStream.Next() {
		event := res.EventStream.Value()
		if event == nil {
			continue
		}

		chunk := event.Data

		if chunk.Error != nil {
			log.Fatalf("errore stream %d: %s", chunk.Error.Code, chunk.Error.Message)
		}

		for _, choice := range chunk.Choices {
			part, ok := choice.Delta.Content.GetOrZero()
			if ok {
				fmt.Print(formatter.FormatChunk(part))
				full.WriteString(part)
			}
		}
	}

	fmt.Print(formatter.Close())

	if err := res.EventStream.Err(); err != nil {
		log.Fatal(err)
	}

	return full.String()
}

func moveCursor(builder *strings.Builder, x, y int) {
	_, err := fmt.Fprintf(builder, "\033[%d;%dH", y, x)
	if err != nil {
		return
	}
}

func HomePage() *Page {
	ClearTerminal()
	page := NewPage(core.WrapIn(
		"                    \n"+
			"|   |               \n"+
			"|---|,---.,-.-.,---.\n"+
			"|   ||   || | ||---'\n"+
			"`   '`---'` ' '`---'\n"+
			"                    ", core.Green),
		"Questo è un client gratuito per interagire con AI senza consumi di ram e cpu alti\n"+
			"Nessun lag dovuti all'app\n"+
			"ESC per uscire\n"+
			"1 - Home\n"+
			"2 - Chat\n"+
			"3 - Chat coding LOW\n"+
			"4 - Chat coding HIGH\n"+
			"5 - Usage and info\n"+
			"6 - Ai Config\n"+
			"7 - Credits", false)
	page.Update()

	return page
}

func Run() {
	currentPage = HomePage()
	handleShortcuts()
}

func handleShortcuts() {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			panic(err)
		}
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case key == keyboard.KeyCtrlZ:
		case char == 'r':
			currentPage.Update()

		case char == '1':
			currentPage = HomePage()
		case char == '2':
			currentPage = ChatPage(currentPage, core.DefaultModel())
		case char == '3':
			currentPage = ChatPage(currentPage, core.CodingLowModel())
		case char == '4':
			currentPage = ChatPage(currentPage, core.CodingHighModel())
		case char == '5':
			currentPage = UsagePage(currentPage)
		case char == '6':
			currentPage = ConfigPage(currentPage)
		case char == '7':
			currentPage = CreditPage(currentPage)
		case key == keyboard.KeyEsc:
			fmt.Println("Addio")
			return
		}

		if key == keyboard.KeyEsc {
			return
		}

	}
}
