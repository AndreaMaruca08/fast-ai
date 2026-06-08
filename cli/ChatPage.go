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

func ChatPage(page *Page, model core.Model) *Page {
	ClearTerminal()

	page = NewPage("Chat", "Modello attuale: "+model.Name, true)
	currentPage = page
	chat = &core.Chat{Model: model}
	page.Update()

	handleCommands(page)

	return page
}
func handleCommands(page *Page) {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}
		switch {
		case key == keyboard.KeyEnter:
			resp, err := getSingleInput(">")
			if err != nil {
				log.Fatal(err)
			}
			if resp == "" {
				return
			}
			send(page, resp)

		case key == keyboard.KeyBackspace:
			return
		}
	}
}
func send(page *Page, message string) {
	if chat == nil {
		chat = &core.Chat{Model: core.DefaultModel()}
	}

	chat.CreateUserMessage(message)
	page.AddMessage(message, true)
	page.Update()

	resp := chat.Send(streamer)
	page.AddMessage(resp, false)
	page.Update()
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

	if scanner.Text() == "" {
		HomePage()
		keyboard.Open()
		return "", nil
	}

	return scanner.Text(), nil
}
