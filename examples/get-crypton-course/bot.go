package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

func (app *botApp) handleRateRequest(sendTo string, dest destinationType) error {
	rate, err := getCryptonRate(app.coingecko)
	if err != nil {
		log.Println("error:", err.Error())
		return nil
	}

	if rate == 0 {
		return fmt.Errorf("at the moment Crypton rate is unknown to me, try again later")
	}

	msg := fmt.Sprintf("Crypton rate is 1 CRP = %v USD", rate)
	return app.sendMessage(sendTo, msg, dest)
}

func (app *botApp) sendMessage(sendTo string, msg string, dest destinationType) error {
	var err error
	switch dest {
	case destTypeChannel:
		_, err = app.chatbot.GetClient().SendInstantMessage(sendTo, msg)
	case destTypeContact:
		_, err = app.chatbot.GetClient().SendChannelMessage(sendTo, msg)
	}
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}

func (app *botApp) OnContactMessage(m structs.InstantMessage) {
	if !strings.Contains(m.Text, botCommand) {
		return
	}

	if err := app.handleRateRequest(m.Pubkey, destTypeContact); err != nil {
		log.Println("error: ", err.Error())

		if err := app.sendMessage(
			m.Pubkey, defaultErrorMessage, destTypeContact,
		); err != nil {
			log.Println("failed to send contact message:", err.Error())
		}
	}
}

func (app *botApp) OnChannelMessage(m structs.WsChannelMessage) {
	if !strings.Contains(m.Text, botCommand) {
		return
	}

	if err := app.handleRateRequest(m.ChannelID, destTypeChannel); err != nil {
		log.Println("error: ", err.Error())

		if err := app.sendMessage(
			m.ChannelID, defaultErrorMessage, destTypeChannel,
		); err != nil {
			log.Println("failed to send channel message:", err.Error())
		}
	}
}

func (app *botApp) OnPrivateChannelMessage(m structs.WsChannelMessage) {
	fmt.Printf("[PRIVATE] %s: %s\n", m.Nick, m.Text)
}

func (app *botApp) OnWelcomeMessage(userPubkey string) string {
	return fmt.Sprintf(
		"Hello! I can show Crypton in private messages "+
			"and in a general chat, to which the administrator will add me.\n\n"+
			"Try entering: %s", botCommand,
	)
}

func onError(err error) {
	if err == nil {
		return
	}

	log.Println("error:", err.Error())
}
