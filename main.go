package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "strings"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    resp,err := http.Get("https://meme-api.herokuapp.com/gimme")
    if err != nil {
        fmt.Fprintln(os.Stderr, "HTTP request to meme API failed")
    }
    defer resp.Body.Close()
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    url := fmt.Sprintf("%v", result["url"])

    resp2,err := http.Get("https://random-word-api.herokuapp.com/word?number=1&swear=0")
    if err != nil {
        fmt.Fprintln(os.Stderr, "HTTP request to random word API failed")
    }
    defer resp2.Body.Close()
    body,err := ioutil.ReadAll(resp2.Body)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to parse body")
    }

    word := string(body)
    word = strings.ReplaceAll(word, "[", "")
    word = strings.ReplaceAll(word, "]", "")
    word = strings.ReplaceAll(word, "\"", "")

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!meme" {
		s.ChannelMessageSend(m.ChannelID, url)
	}

	if m.Content == "Do you like The Quarter Game" {
		s.ChannelMessageSend(m.ChannelID, "YESS! The Quarter Game is better than " + word + "!!!")
	}
}
