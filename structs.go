package uchatbot

import utopiago "github.com/Sagleft/utopialib-go"

type wsHandler func(event utopiago.WsEvent)

type ChatBot struct {
	data       ChatBotData
	wsHandlers map[string]wsHandler
}

type ChatBotCallbacks struct {
	OnContactMessage        func(utopiago.InstantMessage)
	OnChannelMessage        func(utopiago.ChannelMessage)
	OnPrivateChannelMessage func(utopiago.ChannelMessage)
}

type ChatBotData struct {
	// required
	Client    *utopiago.UtopiaClient `json:"client"`
	Chats     []Chat                 `json:"chats"` // channel ids
	Callbacks ChatBotCallbacks

	// optional
	UseErrorCallback bool `json:"useErrorCallback"`
	EnableWsSSL      bool `json:"enableSSL"` // SSL for websocket connection
	ErrorCallback    func(err error)
}

type Chat struct {
	// required
	ID string `json:"id"`

	// optional
	Password string `json:"password"`
}
