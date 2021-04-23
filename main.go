package main

import (
	"fmt"
	"os"
	"github.com/bwmarrin/discordgo"
)

func main() {
	// open config file
	configJSON, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file.")
		return
	}

	// read the configuration settings from the config file.
	fmt.Println("Loading config values from \"config.json\"...")
	byteValue, _ := ioutil.ReadAll(configJSON)
	var configValues map[string]interface{}
	json.Unmarshal([]byte(byteValue), &configValues)
	auth_token := configValues["database"].(string)
	fmt.Println("\033[32mSuccessfully loaded config values!\033[0m")

	// start running the bot
	discord, err := discordgo.New("GitHub Repo Notifier " + auth_token)
	if err != nil {
		fmt.Println("Error starting bot.")
		return
	}
}