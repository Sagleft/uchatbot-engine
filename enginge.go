package uchatbot

import (
	"errors"
	"log"

	utopiago "github.com/Sagleft/utopialib-go"
)

// NewChatBot - create new chatbot and connect to Utopia.
// the bot will try to join the list of the specified chats and subscribe to messages
func NewChatBot(data ChatBotData) (*ChatBot, error) {
	// check data
	if data.Client.WsPort == 0 {
		return nil, errors.New("ws port is not set")
	}

	// create bot
	cb := &ChatBot{
		data:       data,
		wsHandlers: make(map[string]wsHandler),
	}

	return cb, checkErrors(
		cb.checkConnection,
		cb.joinChannels,
		cb.initHandlers,
		cb.subscribe,
	)
}

func (c *ChatBot) checkConnection() error {
	if !c.data.Client.CheckClientConnection() {
		return errors.New("failed to connect to Utopia Client")
	}
	return nil
}

func (c *ChatBot) joinChannels() error {
	for _, chat := range c.data.Chats {
		isJoined, err := c.data.Client.JoinChannel(chat.ID, chat.Password)
		if err != nil {
			return err
		}

		if !isJoined {
			return errors.New("failed to join in " + chat.ID)
		}
	}
	return nil
}

func (c *ChatBot) subscribe() error {
	err := c.data.Client.SetWebSocketState(utopiago.SetWsStateTask{
		Enabled:       true,
		Port:          c.data.Client.WsPort,
		EnableSSL:     c.data.EnableWsSSL,
		Notifications: "all",
	})
	if err != nil {
		return err
	}

	return c.data.Client.WsSubscribe(utopiago.WsSubscribeTask{
		OnConnected: c.onConnected,
		Callback:    c.onMessage,
		ErrCallback: c.onError,
		Port:        c.data.Client.WsPort,
	})
}

func (c *ChatBot) onConnected() {}

func (c *ChatBot) onError(err error) {
	if err == nil {
		return
	}

	if c.data.UseErrorCallback {
		c.data.ErrorCallback(err)
	} else {
		log.Println("ERROR: " + err.Error())
	}
}
