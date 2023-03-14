package uchatbot

import (
	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go/v2"
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
	"github.com/Sagleft/utopialib-go/v2/pkg/websocket"
	"github.com/beefsack/go-rate"
)

type wsHandler func(event websocket.WsEvent)

type ChatBot struct {
	client utopiago.Client

	data         ChatBotData
	wsHandlers   map[string]wsHandler
	queues       eventBuffers
	rateLimiters botRateLimiters
}

type botRateLimiters struct {
	InstantMessage        rateLimiter
	ChannelPrivateMessage rateLimiter
}

type rateLimiter struct {
	L       *rate.RateLimiter
	Enabled bool
}

type eventBuffers struct {
	// receivers
	Auth                *swissknife.ChannelWorker
	Contact             *swissknife.ChannelWorker
	ChannelLobby        *swissknife.ChannelWorker
	PrivateChannelLobby *swissknife.ChannelWorker

	// senders
	InstantMessages           *swissknife.ChannelWorker
	SendPrivateChannelMessage *swissknife.ChannelWorker
}

type ChatBotCallbacks struct {
	// required
	OnContactMessage        func(structs.InstantMessage)
	OnChannelMessage        func(structs.WsChannelMessage)
	OnPrivateChannelMessage func(structs.WsChannelMessage)

	// optional
	WelcomeMessage func(userPubkey string) string
}

type ChatBotData struct {
	// required
	Config    utopiago.Config `json:"client"`
	Chats     []Chat          `json:"chats"` // channel ids
	Callbacks ChatBotCallbacks

	// optional
	SkipConnectionCheck bool                 `json:"skipConnCheck"`
	Notifications       string               `json:"notifications"` // by default: all
	UseErrorCallback    bool                 `json:"useErrorCallback"`
	EnableWsSSL         bool                 `json:"enableSSL"` // SSL for websocket connection
	BuffersCapacity     EventBuffersCapacity `json:"buffersCapacity"`
	RateLimiters        EventBufferLimiters  `json:"rateLimiters"`

	ErrorCallback func(err error) `json:"-"`
}

// for limit max events per second
type EventBufferLimiters struct {
	InstantMessages        int `json:"instantMessages"`
	ChannelPrivateMessages int `json:"channelPrivateMessages"`
}

type EventBuffersCapacity struct {
	Auth                  int `json:"auth"`
	ContactMessage        int `json:"contactMessage"`
	ChannelMessage        int `json:"channelMessage"`
	PrivateChannelMessage int `json:"privateChannelMessage"`
	InstantMessages       int `json:"instantMessages"`
}

type Chat struct {
	// required
	ID string `json:"id"`

	// optional
	Password string `json:"password"`
}
