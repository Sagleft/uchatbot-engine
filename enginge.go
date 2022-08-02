package uchatbot

import (
	"errors"

	utopiago "github.com/Sagleft/utopialib-go"
)

type ChatBot struct {
	data ChatBotData
}

type ChatBotData struct {
	// required
	Client *utopiago.UtopiaClient `json:"client"`
	Chats  []Chat                 `json:"chats"` // channel ids

	// optional
	UseErrorCallback bool
	ErrorCallback    func(err error)
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
	// check data
	if data.Client.WsPort == 0 {
		return nil, errors.New("ws port is not set")
	}

	// create bot
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

	// subscribe to events
	err := cb.data.Client.WsSubscribe(utopiago.WsSubscribeTask{
		OnConnected: cb.onConnected,
		Callback:    cb.onMessage,
		ErrCallback: cb.onError,
		Port:        data.Client.WsPort,
	})
	if err != nil {
		return cb, err
	}

	return cb, nil
}

func (c *ChatBot) onConnected() {}

func (c *ChatBot) onMessage(ws utopiago.WsEvent) {
	// TODO
}

func (c *ChatBot) onError(err error) {
	if c.data.UseErrorCallback {
		c.data.ErrorCallback(err)
	} else {
		// TODO
	}
}
