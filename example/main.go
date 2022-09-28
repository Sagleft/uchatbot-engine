package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go"
)

func main() {
	APIHost := os.Getenv("HOST")
	APIToken := os.Getenv("TOKEN")

	if APIHost == "" {
		log.Fatalln("API host is not set")
	}

	if APIToken == "" {
		log.Fatalln("API token is not set")
	}

	_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Client: &utopiago.UtopiaClient{
			Protocol: "http",
			Host:     APIHost,
			Token:    APIToken,
			Port:     22800,
			WsPort:   25000,
		},
		Chats: []uchatbot.Chat{
			{ID: "D53B4431FD604E2F0261792444797AA4"},
			{ID: "A59D8B62E1A59049564A4B0F8B457D45"},
		},
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        OnContactMessage,
			OnChannelMessage:        OnChannelMessage,
			OnPrivateChannelMessage: OnPrivateChannelMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    onError,
	})
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("connected")
	}

	wait()
}

func wait() {
	ch := make(chan struct{})
	// run in background
	<-ch
}

func onError(err error) {
	log.Println("ERROR: " + err.Error())
}

func OnContactMessage(m utopiago.InstantMessage) {
	fmt.Println("[CONTACT] " + m.Nick + ": " + m.Text)
}

func OnChannelMessage(m utopiago.WsChannelMessage) {
	fmt.Println("[CHANNEL] " + m.Nick + ": " + m.Text)
}

func OnPrivateChannelMessage(m utopiago.WsChannelMessage) {
	fmt.Println("[PRIVATE] " + m.Nick + ": " + m.Text)
}
