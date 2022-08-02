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
	Chats  []Chat                 `json:"chats"` // channel ids
}

type Chat struct {
	// required
	ID string `json:"id"`

	// optional
	Password string `json:"password"`
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

	// join to channels (chats)
	for _, chat := range cb.data.Chats {
		isJoined, err := cb.data.Client.JoinChannel(chat.ID, chat.Password)
		if err != nil {
			return cb, err
		}

		if !isJoined {
			return cb, errors.New("failed to join in " + chat.ID)
		}
	}

	return cb, nil
}
