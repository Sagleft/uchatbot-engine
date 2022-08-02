package uchatbot

import (
	utopiago "github.com/Sagleft/utopialib-go"
)

func (c *ChatBot) initHandlers() error {
	c.wsHandlers = map[string]wsHandler{
		"newInstantMessage":        c.onContactMessage,
		"newChannelMessage":        c.onChannelMessage,
		"newPrivateChannelMessage": c.onPrivateChannelMessage,
	}
	return nil
}

func (c *ChatBot) onMessage(event utopiago.WsEvent) {
	handler, isEventTypeKown := c.wsHandlers[event.Type]
	if !isEventTypeKown {
		return
	}

	handler(event)
}

/*{
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
	message, err := event.GetInstantMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnContactMessage(message)
}

/*{
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
	message, err := event.GetChannelMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnChannelMessage(message)
}

/*{
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
	message, err := event.GetChannelMessage()
	if err != nil {
		c.onError(err)
		return
	}

	c.data.Callbacks.OnPrivateChannelMessage(message)
}
