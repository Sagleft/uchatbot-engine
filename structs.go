package uchatbot

import (
	swissknife "github.com/Sagleft/swiss-knife"
	utopiago "github.com/Sagleft/utopialib-go"
)

type wsHandler func(event utopiago.WsEvent)

type ChatBot struct {
	data       ChatBotData
	wsHandlers map[string]wsHandler
	queues     eventBuffers
}

type eventBuffers struct {
	Auth                *swissknife.ChannelWorker
	Contact             *swissknife.ChannelWorker
	ChannelLobby        *swissknife.ChannelWorker
	PrivateChannelLobby *swissknife.ChannelWorker
}

type ChatBotCallbacks struct {
	OnContactMessage        func(utopiago.InstantMessage)
	OnChannelMessage        func(utopiago.WsChannelMessage)
	OnPrivateChannelMessage func(utopiago.WsChannelMessage)
}

type ChatBotData struct {
	// required
	Client    *utopiago.UtopiaClient `json:"client"`
	Chats     []Chat                 `json:"chats"` // channel ids
	Callbacks ChatBotCallbacks

	// optional
	UseErrorCallback bool                 `json:"useErrorCallback"`
	EnableWsSSL      bool                 `json:"enableSSL"` // SSL for websocket connection
	ErrorCallback    func(err error)      `json:"-"`
	BuffersCapacity  EventBuffersCapacity `json:"buffersCapacity"`
}

type EventBuffersCapacity struct {
	Auth                  int `json:"auth"`
	ContactMessage        int `json:"contactMessage"`
	ChannelMessage        int `json:"channelMessage"`
	PrivateChannelMessage int `json:"privateChannelMessage"`
}

type Chat struct {
	// required
	ID string `json:"id"`

	// optional
	Password string `json:"password"`
}
