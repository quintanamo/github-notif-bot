package main

import (
	"fmt"
	"os"
	"os/signal"
	"encoding/json"
    "io/ioutil"
	"syscall"
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
	auth_token := configValues["bot-auth-token"].(string)
	fmt.Println("\033[32mSuccessfully loaded config values!\033[0m")

	// start running the bot
	discord, err := discordgo.New("Bot " + auth_token)
	if err != nil {
		fmt.Println("Error initializing bot.")
		return
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	//discord.ChannelMessageSend("835167014720897064", "Hello!")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the Discord session.
	discord.Close()
}