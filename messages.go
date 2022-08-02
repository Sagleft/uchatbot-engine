package uchatbot

import (
	"log"

	utopiago "github.com/Sagleft/utopialib-go"
)

func (c *ChatBot) onMessage(event utopiago.WsEvent) {
	log.Println(event.Type)
}
