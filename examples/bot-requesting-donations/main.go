package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Sagleft/uchatbot-engine"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"

	"github.com/Sagleft/cronlib"
)

const APIToken = "your-utopia-api-token"
const utopiaHost = "127.0.0.1"
const utopiaPort = 20000
const utopiaWsPort = 25000
const sendDonateRequestTimeout = time.Hour * 24
const sendDonateRequestAtStart = true
const useUCode = false
const uCodeComment = ""

const findBotMessageInLast = 10

var chats = []uchatbot.Chat{
	{ID: "D53B4431FD604E2F0261792444797AA4"},
}

type botApp struct {
	chatbot *uchatbot.ChatBot
}

func main() {
	var app = botApp{}
	var err error

	cronlib.NewCronHandler(app.onTimer, sendDonateRequestTimeout).
		Run(sendDonateRequestAtStart)

	app.chatbot, err = uchatbot.NewChatBot(uchatbot.ChatBotData{
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

func (a *botApp) onTimer() {
	// let's use first chat (channel) ID
	channelID := chats[0].ID

	// create donate service
	srv := a.chatbot.RequestDonate(channelID).
		EnableUCode(useUCode, uCodeComment)

	_, donateMessage, err := srv.GetDonateMessage()
	if err != nil {
		log.Println(err)
		return
	}

	// just in case, we’ll check if there was a duplicate of this message
	// in the latest chat messages so that the bot doesn’t spam
	lastMessages, err := a.chatbot.GetClient().GetChannelMessages(
		channelID, 0, findBotMessageInLast,
	)
	if err != nil {
		log.Println(err)
		return
	}

	for _, channelMessage := range lastMessages {
		if strings.Contains(channelMessage.Text, donateMessage) {
			return // skip. the bot recently sent a message asking for donations
		}
	}

	// send request to channel
	if err := srv.Do(); err != nil {
		log.Println(err)
		return
	}
}
