package uchatbot

import (
	utopiago "github.com/Sagleft/utopialib-go"
)

/*
                _                    _        _
  __      _____| |__  ___  ___   ___| | _____| |_
  \ \ /\ / / _ \ '_ \/ __|/ _ \ / __| |/ / _ \ __|
   \ V  V /  __/ |_) \__ \ (_) | (__|   <  __/ |_
    \_/\_/ \___|_.__/|___/\___/ \___|_|\_\___|\__|

*/

func (c *ChatBot) onMessage(event utopiago.WsEvent) {
	handler, isEventTypeKown := c.wsHandlers[event.Type]
	if !isEventTypeKown {
		return
	}

	handler(event)
}

func (c *ChatBot) initHandlers() error {
	c.wsHandlers = map[string]wsHandler{
		"newAuthorization":         c.onNewAuth,
		"newInstantMessage":        c.onContactMessage,
		"newChannelMessage":        c.onChannelMessage,
		"newPrivateChannelMessage": c.onPrivateChannelMessage,
	}
	return nil
}

/*
               _   _
    __ _ _   _| |_| |__
   / _` | | | | __| '_ \
  | (_| | |_| | |_| | | |
   \__,_|\__,_|\__|_| |_|

*/

func (c *ChatBot) onNewAuth(event utopiago.WsEvent) {
	c.queues.Auth.AddEvent(event)
}

func (c *ChatBot) handleAuthEvent(e interface{}) {
	event, err := c.convertEventInterface(e, "auth")
	if err != nil {
		c.onError(err)
		return
	}

	// get pubkey
	userPubkey, err := event.GetString("pk")
	if err != nil {
		c.onError(err)
		return
	}

	// approve auth
	_, err = c.data.Client.AcceptAuthRequest(userPubkey, "")
	if err != nil {
		c.onError(err)
		return
	}

	// send welcome message
	if c.data.Callbacks.WelcomeMessage == nil {
		return // callback is not set
	}

	msgText := c.data.Callbacks.WelcomeMessage(userPubkey)
	c.SendContactMessage(userPubkey, msgText)
}

/*
                   _             _
    ___ ___  _ __ | |_ __ _  ___| |_    _ __ ___  ___  __ _
   / __/ _ \| '_ \| __/ _` |/ __| __|  | '_ ` _ \/ __|/ _` |
  | (_| (_) | | | | || (_| | (__| |_   | | | | | \__ \ (_| |
   \___\___/|_| |_|\__\__,_|\___|\__|  |_| |_| |_|___/\__, |
                                                      |___/
{
	"type": "newInstantMessage",
	"data": {
		"dateTime": "2022-08-02T18:54:14.437Z",
		"file": null,
		"id": 609,
		"isIncoming": true,
		"messageType": 1,
		"metaData": null,
		"nick": "Tester",
		"pk": "150AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804",
		"readDateTime": null,
		"receivedDateTime": "2022-08-02T18:54:14.437Z",
		"text": "test"
	}
}*/

func (c *ChatBot) onContactMessage(event utopiago.WsEvent) {
	c.queues.Contact.AddEvent(event)
}

func (c *ChatBot) handleContactMessage(e interface{}) {
	event, err := c.convertEventInterface(e, "contact message")
	if err != nil {
		c.onError(err)
		return
	}

	message, err := event.GetInstantMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnContactMessage(message)
}

/*
        _                            _
    ___| |__   __ _ _ __  _ __   ___| |   _ __ ___  ___  __ _
   / __| '_ \ / _` | '_ \| '_ \ / _ \ |  | '_ ` _ \/ __|/ _` |
  | (__| | | | (_| | | | | | | |  __/ |  | | | | | \__ \ (_| |
   \___|_| |_|\__,_|_| |_|_| |_|\___|_|  |_| |_| |_|___/\__, |
                                                        |___/

{
	"type": "newChannelMessage",
	"data": {
		"channel": "bot playground",
		"channelid": "A59D8B62E1A59049564A4B0F8B457D45",
		"dateTime": "2022-08-02T19:00:13.188Z",
		"hashedPk": "2171431BFE9F9770A9A2794A2A203E01",
		"isIncoming": true,
		"messageType": 1,
		"metaData": null,
		"nick": "Tester",
		"pk": "F50AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804",
		"text": "test",
		"topicId": "4118503780692390786"
	}
}*/

func (c *ChatBot) onChannelMessage(event utopiago.WsEvent) {
	c.queues.ChannelLobby.AddEvent(event)
}

func (c *ChatBot) handleChannelLobbyMessage(e interface{}) {
	event, err := c.convertEventInterface(e, "channel message")
	if err != nil {
		c.onError(err)
		return
	}

	message, err := event.GetChannelMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnChannelMessage(message)
}

/*
        _                            _               _
    ___| |__   __ _ _ __  _ __   ___| |   _ __  _ __(_)_   __   _ __ ___  ___  __ _
   / __| '_ \ / _` | '_ \| '_ \ / _ \ |  | '_ \| '__| \ \ / /  | '_ ` _ \/ __|/ _` |
  | (__| | | | (_| | | | | | | |  __/ |  | |_) | |  | |\ V /   | | | | | \__ \ (_| |
   \___|_| |_|\__,_|_| |_|_| |_|\___|_|  | .__/|_|  |_| \_/    |_| |_| |_|___/\__, |
                                         |_|                                  |___/

{
	"type": "newPrivateChannelMessage",
	"data": {
		"channel": "bot playground",
		"channelid": "A59D8B62E1A59049564A4B0F8B457D45",
		"dateTime": "2022-08-02T19:04:39.536Z",
		"hashedPk": "2171431BFE9F9770A9A2794A2A203E01",
		"isIncoming": true,
		"messageType": 1,
		"metaData": null,
		"nick": "Tester",
		"pk": "F50AF5410B1F3F4297043F0E046F205BCBAA76BEC70E936EB0F3AB94BF316804",
		"text": "test",
		"topicId": "2679238831019856860"
	}
}*/

func (c *ChatBot) onPrivateChannelMessage(event utopiago.WsEvent) {
	c.queues.PrivateChannelLobby.AddEvent(event)
}

func (c *ChatBot) handlePrivateChannelLobbyMessage(e interface{}) {
	event, err := c.convertEventInterface(e, "private channel message")
	if err != nil {
		c.onError(err)
		return
	}

	message, err := event.GetChannelMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnPrivateChannelMessage(message)
}
