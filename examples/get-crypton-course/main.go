package main

import (
	"log"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
)

func main() {
	if err := runBot(); err != nil {
		log.Println(err)
	}
}

func runBot() error {
	var app = botApp{
		coingecko: crateCoingeckoClient(),
	}
	var err error

	app.chatbot, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
		Config: utopiago.Config{
			Host:   utopiaHost,
			Token:  APIToken,
			Port:   utopiaPort,
			WsPort: utopiaWsPort,
		},
		Chats: chats,
		Callbacks: uchatbot.ChatBotCallbacks{
			OnContactMessage:        app.OnContactMessage,
			OnChannelMessage:        app.OnChannelMessage,
			OnPrivateChannelMessage: app.OnPrivateChannelMessage,

			WelcomeMessage: app.OnWelcomeMessage,
		},
		UseErrorCallback: true,
		ErrorCallback:    onError,
	})
	return err
}
