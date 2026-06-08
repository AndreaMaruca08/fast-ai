package main

import (
	"fast_ai_client/cli"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	cli.Run()
}

func loadEnv() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	exeDir := filepath.Dir(exePath)
	envPath := filepath.Join(exeDir, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("nessun .env trovato in %s", envPath)
	}
}
