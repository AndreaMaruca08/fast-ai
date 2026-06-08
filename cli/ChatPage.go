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

	page = NewPage(core.WrapIn(
		"  ░██████  ░██                      ░██    \n"+
			" ░██   ░██ ░██                      ░██    \n"+
			"░██        ░████████   ░██████   ░████████ \n"+
			"░██        ░██    ░██       ░██     ░██    \n"+
			"░██        ░██    ░██  ░███████     ░██    \n"+
			" ░██   ░██ ░██    ░██ ░██   ░██     ░██    \n"+
			"  ░██████  ░██    ░██  ░█████░██     ░████ \n"+
			"                                           \n", core.Green),
		"Modello attuale: "+core.WrapIn(model.Name, core.Red)+" | "+core.WrapIn(model.Type, core.Blue),
		true,
	)
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

	resp, err := chat.Send(streamer)
	if err != nil {
		fmt.Println(err)
		return
	}
	page.AddMessage(resp, false)
	page.AddMessage("\n--------------------------", false)
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

	if scanner.Text() == "" || scanner.Text() == "home" || scanner.Text() == "exit" {
		HomePage()
		err := keyboard.Open()
		if err != nil {
			return "", err
		}
		return "", nil
	}

	return scanner.Text(), nil
}
