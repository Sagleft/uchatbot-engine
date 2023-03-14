package uchatbot

import (
	"errors"
	"log"
	"time"

	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	uErrors "github.com/Sagleft/utopialib-go/v2/pkg/errors"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/Sagleft/utopialib-go/v2/pkg/websocket"
	"github.com/beefsack/go-rate"
)

const (
	defaultBufferCapacity = 150 // events
	wsReconnectTimeout    = 5 * time.Second
)

// NewChatBot - create new chatbot and connect to Utopia.
// the bot will try to join the list of the specified chats and subscribe to messages
func NewChatBot(data ChatBotData) (*ChatBot, error) {
	// check data
	if data.Config.WsPort == 0 {
		return nil, errors.New("ws port is not set")
	}

	// create bot
	cb := &ChatBot{
		data:       data,
		client:     utopiago.NewUtopiaClient(data.Config),
		wsHandlers: make(map[string]wsHandler),
	}

	return cb, checkErrors(
		cb.checkConnection,
		cb.joinChannels,
		cb.setupMessageQueues,
		cb.initHandlers,
		cb.subscribe,
	)
}

func (c *ChatBot) checkConnection() error {
	if c.data.SkipConnectionCheck {
		return nil
	}

	if !c.client.CheckClientConnection() {
		return errors.New("failed to connect to Utopia Client at `" +
			c.data.Config.Host + "`")
	}
	return nil
}

func (c *ChatBot) joinChannels() error {
	for _, chat := range c.data.Chats {
		isJoined, err := c.client.JoinChannel(chat.ID, chat.Password)
		if err != nil {
			return err
		}

		if !isJoined {
			return errors.New("failed to join in " + chat.ID)
		}
	}
	return nil
}

func (c *ChatBot) setupMessageQueues() error {
	// RECEIVERS
	c.queues.Auth = swissknife.NewChannelWorker(
		c.handleAuthEvent,
		ternaryInt(
			c.data.BuffersCapacity.Auth == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.Auth,
		),
	)
	go c.queues.Auth.Start()

	c.queues.Contact = swissknife.NewChannelWorker(
		c.handleContactMessage,
		ternaryInt(
			c.data.BuffersCapacity.ContactMessage == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.ContactMessage,
		),
	)
	go c.queues.Contact.Start()

	c.queues.ChannelLobby = swissknife.NewChannelWorker(
		c.handleChannelLobbyMessage,
		ternaryInt(
			c.data.BuffersCapacity.ChannelMessage == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.ChannelMessage,
		),
	)
	go c.queues.ChannelLobby.Start()

	c.queues.PrivateChannelLobby = swissknife.NewChannelWorker(
		c.handlePrivateChannelLobbyMessage,
		ternaryInt(
			c.data.BuffersCapacity.PrivateChannelMessage == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.PrivateChannelMessage,
		),
	)
	go c.queues.PrivateChannelLobby.Start()

	// SENDERS
	c.queues.InstantMessages = swissknife.NewChannelWorker(
		c.handleSendInstantMessageTask,
		ternaryInt(
			c.data.BuffersCapacity.InstantMessages == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.InstantMessages,
		),
	)
	if c.data.RateLimiters.InstantMessages > 0 {
		c.rateLimiters.InstantMessage = rateLimiter{
			L:       rate.New(c.data.RateLimiters.InstantMessages, time.Second),
			Enabled: true,
		}
	}
	go c.queues.InstantMessages.Start()

	c.queues.SendPrivateChannelMessage = swissknife.NewChannelWorker(
		c.handleSendPrivateChannelMessageTask,
		ternaryInt(
			c.data.BuffersCapacity.PrivateChannelMessage == 0,
			defaultBufferCapacity,
			c.data.BuffersCapacity.PrivateChannelMessage,
		),
	)
	if c.data.RateLimiters.ChannelPrivateMessages > 0 {
		c.rateLimiters.InstantMessage = rateLimiter{
			L:       rate.New(c.data.RateLimiters.ChannelPrivateMessages, time.Second),
			Enabled: true,
		}
	}
	go c.queues.SendPrivateChannelMessage.Start()
	return nil
}

func (c *ChatBot) subscribe() error {
	if c.data.DisableEvents {
		return nil
	}

	if c.data.Notifications == "" {
		c.data.Notifications = "all"
	}

	err := c.client.SetWebSocketState(structs.SetWsStateTask{
		Enabled:       true,
		Port:          c.data.Config.WsPort,
		EnableSSL:     c.data.EnableWsSSL,
		Notifications: c.data.Notifications,
	})
	if err != nil {
		return err
	}

	return c.client.WsSubscribe(websocket.WsSubscribeTask{
		OnConnected: c.onConnected,
		Callback:    c.onMessage,
		ErrCallback: c.onWsError,
	})
}

func (c *ChatBot) onConnected() {}

func (c *ChatBot) onWsError(err error) {
	if err == nil {
		return
	}

	if uErrors.CheckErrorConnBroken(err) {
		c.onError(errors.New("websocket connection closed. attempt to reconnect"))

		for {
			if err := c.subscribe(); err == nil {
				return
			}
			time.Sleep(wsReconnectTimeout)
		}
	}

	c.onError(err)
}

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
