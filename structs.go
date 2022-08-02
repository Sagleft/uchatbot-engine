package uchatbot

import utopiago "github.com/Sagleft/utopialib-go"

type wsHandler func(event utopiago.WsEvent)

type ChatBot struct {
	data       ChatBotData
	wsHandlers map[string]wsHandler
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
