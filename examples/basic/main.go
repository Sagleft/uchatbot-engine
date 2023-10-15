package main

import (
	"fmt"
	"log"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

const APIToken = "your-utopia-api-token"
const utopiaHost = "127.0.0.1"
const utopiaPort = 20000
const utopiaWsPort = 25000

var chats = []uchatbot.Chat{
	{ID: "D53B4431FD604E2F0261792444797AA4"},
	{ID: "A59D8B62E1A59049564A4B0F8B457D45"},
}

func main() {
	_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: utopiago.Config{
			Host:   utopiaHost,
			Token:  APIToken,
			Port:   utopiaPort,
			WsPort: utopiaWsPort,
		},
		Chats: chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        OnContactMessage,
			OnChannelMessage:        OnChannelMessage,
			OnPrivateChannelMessage: OnPrivateChannelMessage,

			WelcomeMessage: OnWelcomeMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    onError,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func OnContactMessage(m structs.InstantMessage) {
	fmt.Printf("[CONTACT] %s: %s\n", m.Nick, m.Text)
}

func OnChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[CHANNEL] %s: %s\n", m.Nick, m.Text)
}

func OnPrivateChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[PRIVATE] %s: %s\n", m.Nick, m.Text)
}

func OnWelcomeMessage(userPubkey string) string {
	return fmt.Sprintf("Hello! Your pubkey is %s", userPubkey)
}

func onError(err error) {
	if err == nil {
		return
	}

	log.Println("error:", err.Error())
}
