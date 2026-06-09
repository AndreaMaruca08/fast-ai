package cli

import (
	"bufio"
	"fast_ai_client/core"
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

var chat *core.Chat
var streamer = &Streamer{}
var currentHistory *core.ChatHistory
var current int
var histories []core.ChatHistory

func ChatPage(page *Page, model core.Model) *Page {
	ClearTerminal()
	page = NewPage(core.WrapIn(
		"  ░██████  ░██                      ░██    \n"+
			" ░██   ░██ ░██                      ░██    \n"+
			"░██        ░████████   ░██████   ░████████ \n"+
			"░██        ░██    ░██       ░██     ░██    \n"+
			"░██        ░██    ░██  ░███████     ░██    \n"+
			" ░██   ░██ ░██    ░██ ░██   ░██     ░██    \n"+
			"  ░██████  ░██    ░██  ░█████░██     ░████ \n"+
			"                                           \n", core.Green),
		"Cambia modello con numeri 1, 2, 3\n"+
			"Premi invio con messaggio vuoto per tornare alla home\n"+
			"Premi 'c' per caricare la chat corrente o invio per iniziare una nuova\n"+
			"Modello attuale: "+core.WrapIn(model.Name, core.Red)+" | "+core.WrapIn(model.Type, core.Blue),
		true,
	)
	displaySelection(page)

	if chat == nil {
		chat = &core.Chat{Model: model}
	} else {
		chat.Model = model
	}
	page.Update()

	return page
}
func displaySelection(page *Page) {
	histories = core.GetHistories()
	if len(histories) > 0 {
		if current >= len(histories) {
			current = 0
		}
		for i, history := range histories {
			msg := history.String()
			if i == current {
				msg += " <"
			}
			page.AddMessage(msg, false)
		}
	}
}

func handleChatKey(char rune, key keyboard.Key) bool {
	switch {
	case key == keyboard.KeyArrowDown:
		if len(histories) == 0 {
			return false
		}
		current++
		if current >= len(histories) {
			current = 0
		}
		currentPage = ChatPage(currentPage, chat.Model)
	case char == 'r':
		currentPage.Update()
	case char == '1':
		currentPage = ChatPage(currentPage, core.DefaultModel())
	case char == '2':
		currentPage = ChatPage(currentPage, core.CodingLowModel())
	case char == '3':
		currentPage = ChatPage(currentPage, core.CodingHighModel())
	case char == 'c':
		if len(histories) == 0 || current >= len(histories) {
			return false
		}
		currentHistory = &histories[current]
		chat = &currentHistory.Chat
		currentPage = ChatPage(currentPage, chat.Model)
		currentPage.AddMessage(currentHistory.Chat.String(), false)
		currentPage.Update()
	case key == keyboard.KeyEsc:
		currentPage = HomePage()
	case key == keyboard.KeyEnter:
		resp, err := getSingleInput("\n > ")
		if err != nil {
			log.Fatal(err)
		}
		if resp == "" || resp == "home" || resp == "exit" {
			currentPage = HomePage()
			return false
		}
		send(currentPage, resp)
	case key == keyboard.KeyBackspace:
		currentPage = HomePage()
	}

	return false
}
func send(page *Page, message string) {
	if chat == nil {
		chat = &core.Chat{Model: core.DefaultModel()}
	}

	chat.CreateUserMessage(message)
	page.AddMessage(message, true)
	page.Update()

	resp, err := chat.Send(streamer)
	if err != nil {
		fmt.Println(err)
		return
	}
	page.AddMessage(resp, false)
	page.AddMessage("\n____________________________________________________________________", false)
	page.Update()
	saveCurrentChat()
}

func saveCurrentChat() {
	if currentHistory == nil {
		currentHistory = core.NewChatHistory(*chat)
	}

	currentHistory.Chat = *chat

	for i, history := range histories {
		if history.Title == currentHistory.Title {
			histories[i] = *currentHistory
			core.SaveHistory(histories)
			return
		}
	}

	histories = append(histories, *currentHistory)
	current = len(histories) - 1
	core.SaveHistory(histories)
}

func getSingleInput(prompt string) (string, error) {
	if err := keyboard.Close(); err != nil {
		return "", err
	}

	defer func() {
		if err := keyboard.Open(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}
