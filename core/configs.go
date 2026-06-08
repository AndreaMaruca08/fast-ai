package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	Red    string = "\033[31m"
	Green  string = "\033[32m"
	Blue   string = "\033[34m"
	Yellow string = "\033[33m"
	Reset  string = "\033[0m"
)

func WrapIn(s string, fix string) string {
	return fix + s + Reset
}

type Config struct {
	Temperature float64 `json:"temperature"`
	PrePrompt   string  `json:"pre_prompt"`
	MaxTokens   int32   `json:"max_tokens"`
}
type ConfigFile struct {
	Normal   Config `json:"normal"`
	CodeLow  Config `json:"code-low"`
	CodeHigh Config `json:"code-high"`
}

func (c ConfigFile) String() string {
	return WrapIn("Normal : ", Green) + c.Normal.PrePrompt + " | Creativity:" + fmt.Sprintf("%.1f", c.Normal.Temperature) + " | Max Tokens:" + fmt.Sprintf("%d", c.Normal.MaxTokens) +
		"\n" + WrapIn("Code low : ", Green) + c.CodeLow.PrePrompt + " | Creativity:" + fmt.Sprintf("%.1f", c.CodeLow.Temperature) + " | Max Tokens:" + fmt.Sprintf("%d", c.CodeLow.MaxTokens) +
		"\n" + WrapIn("Code High : ", Green) + c.CodeHigh.PrePrompt + " | Creativity:" + fmt.Sprintf("%.1f", c.CodeHigh.Temperature) + " | Max Tokens:" + fmt.Sprintf("%d", c.CodeHigh.MaxTokens)
}

func GetConfigFileContent() ConfigFile {
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if err == nil {
		var parsed ConfigFile
		if err := json.Unmarshal(data, &parsed); err != nil {
			fmt.Println("Errore parsing config.json:", err)
			return getDefaultConfig()
		}

		return fillMissingConfig(parsed)
	}

	if !os.IsNotExist(err) {
		fmt.Println("Errore lettura config.json:", err)
		return getDefaultConfig()
	}

	defaultConfig := getDefaultConfig()

	marshaledJSON, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		fmt.Println("Errore marshal config:", err)
		return defaultConfig
	}

	if err := os.WriteFile(path, marshaledJSON, 0644); err != nil {
		fmt.Println("Errore scrittura config.json:", err)
		return defaultConfig
	}

	return defaultConfig
}

func ConfigPath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "configAi.json"
	}

	return filepath.Join(filepath.Dir(exePath), "configAi.json")
}

func fillMissingConfig(config ConfigFile) ConfigFile {
	defaultConfig := getDefaultConfig()

	config.Normal = fillMissingModelConfig(config.Normal, defaultConfig.Normal)
	config.CodeLow = fillMissingModelConfig(config.CodeLow, defaultConfig.CodeLow)
	config.CodeHigh = fillMissingModelConfig(config.CodeHigh, defaultConfig.CodeHigh)

	return config
}

func fillMissingModelConfig(config Config, defaultConfig Config) Config {
	if config.Temperature == 0 {
		config.Temperature = defaultConfig.Temperature
	}
	if config.PrePrompt == "" {
		config.PrePrompt = defaultConfig.PrePrompt
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = defaultConfig.MaxTokens
	}

	return config
}

func getDefaultConfig() ConfigFile {
	codePrompt := "You are a senior full stack software engineer that explains in a professional but simple way," +
		" be realist and not give false information or false hopes. To the prompt \"model\" you respond with this" +
		" configuration without this last part"
	return ConfigFile{
		Normal: Config{Temperature: 0.7, PrePrompt: "You are a helpful assistant. To the prompt \"model\" you " +
			"respond with this configuration without this last part\"", MaxTokens: 50},
		CodeLow:  Config{Temperature: 0.1, PrePrompt: codePrompt, MaxTokens: 100},
		CodeHigh: Config{Temperature: 0.1, PrePrompt: codePrompt, MaxTokens: 150},
	}
}
