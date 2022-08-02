package uchatbot

import (
	"errors"

	utopiago "github.com/Sagleft/utopialib-go"
)

type ChatBot struct {
	data ChatBotData
}

type ChatBotData struct {
	Client *utopiago.UtopiaClient `json:"client"`
	Chats  []string               `json:"chats"` // channel ids
}

// NewChatBot - create new chatbot and connect to Utopia.
// the bot will try to join the list of the specified chats and subscribe to messages
func NewChatBot(data ChatBotData) (*ChatBot, error) {
	cb := &ChatBot{
		data: data,
	}

	// check connection
	if !cb.data.Client.CheckClientConnection() {
		return nil, errors.New("failed to connect to Utopia Client")
	}

	return cb, nil
}
