package main

import (
	"fmt"
	"strconv"
	"os"
	"os/signal"
	"encoding/json"
    "io/ioutil"
	"syscall"
	"net/http"
	"github.com/bwmarrin/discordgo"
)

type Issue struct {
	IssueNum int
	Action string
	Title string
	Actor string
	Url string
}

func WatchRepoIssues(url string, channelID string) {
	for {
		// get issue events from GitHub
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		// read bytes from response and deserialize JSON
		bodyBytes, _ := ioutil.ReadAll(response.Body)
		var events []map[string]interface{}
		json.Unmarshal([]byte(bodyBytes), &events)

		// store result in an Issue
		issue := Issue {}
		issue.IssueNum 	= int(events[0]["issue"].(map[string]interface{})["number"].(float64))
		issue.Action 	= events[0]["event"].(string)
		issue.Title 	= events[0]["issue"].(map[string]interface{})["title"].(string)
		issue.Actor 	= events[0]["actor"].(map[string]interface{})["login"].(string)
		issue.Url	 	= events[0]["issue"].(map[string]interface{})["html_url"].(string)

		// build message string
		message := "Issue #" + strconv.Itoa(issue.IssueNum) + " - " + issue.Title + "\n" + "Action: " + issue.Action + " by " + issue.Actor + "\n" + issue.Url

		// output the message in the specified channel
		SendDiscordMessage(channelID, message)
	}
}

func SendDiscordMessage(channelID string, message string) {
	discord.ChannelMessageSend(channelID, message)
}

var discord *discordgo.Session = new(discordgo.Session)

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
	discord, err = discordgo.New("Bot " + auth_token)
	if err != nil {
		fmt.Println("Error initializing bot.")
		return
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	go WatchRepoIssues("https://api.github.com/repos/YCPRadioTelescope/Radio-Tele-Frontend/issues/events", "835167014720897064")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the Discord session.
	discord.Close()
}