package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"math/rand/v2"
)

type ChatHistory struct {
	Title string `json:"title"`
	Chat  Chat   `json:"chat"`
}

func (h *ChatHistory) String() string {
	return h.Title + " | msg count: " + fmt.Sprint(len(h.Chat.History))
}

func NewChatHistory(chat Chat) *ChatHistory {
	return &ChatHistory{Title: fmt.Sprintf("%d", rand.IntN(1000000)), Chat: chat}
}

func GetHistories() []ChatHistory {
	path := HistoriesPath()

	data, err := os.ReadFile(path)
	if err == nil {
		var parsed []ChatHistory
		if err := json.Unmarshal(data, &parsed); err != nil {
			fmt.Println("Errore parsing histories.json:", err)
			return []ChatHistory{}
		}
		return parsed
	}

	defaultH := []ChatHistory{}

	if !os.IsNotExist(err) {
		fmt.Println("Errore lettura histories.json:", err)
		return defaultH
	}

	marshaledJSON, err := json.MarshalIndent(defaultH, "", "  ")
	if err != nil {
		fmt.Println("Errore marshal histories:", err)
		return defaultH
	}

	if err := os.WriteFile(path, marshaledJSON, 0644); err != nil {
		fmt.Println("Errore scrittura histories.json:", err)
		return defaultH
	}

	return defaultH
}

func SaveHistory(histories []ChatHistory) {
	marshaledJSON, err := json.MarshalIndent(histories, "", "  ")
	if err != nil {
		fmt.Println("Errore marshal histories:", err)
		return
	}

	if err := os.WriteFile(HistoriesPath(), marshaledJSON, 0644); err != nil {
		fmt.Println("Errore scrittura histories.json:", err)
	}
}

func HistoriesPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "histories.json"
	}

	return filepath.Join(filepath.Dir(exePath), "histories.json")
}
