package uchatbot

import (
	"errors"
	"reflect"
)

// SetReadonly - enable or disable channel readonly mode
func (c *ChatBot) SetReadonly(channelID string, readOnly bool) error {
	return c.data.Client.EnableReadOnly(channelID, readOnly)
}

// SendContactMessage - send message to contact.
// it works with queue (buffer).
// returns message ID, error
func (c *ChatBot) SendContactMessage(userPubkey string, msgText string) {
	c.queues.InstantMessages.AddEvent(sendMessageTask{
		UserPubkey:  userPubkey,
		MessageText: msgText,
	})
}

type sendMessageTask struct {
	UserPubkey  string
	MessageText string
}

func (c *ChatBot) handleSendInstantMessageTask(e interface{}) {
	event, isConvertable := e.(sendMessageTask)
	if !isConvertable {
		c.onError(errors.New("failed to convert send message task: " +
			reflect.ValueOf(e).String() + " type received"))
		return
	}

	_, err := c.data.Client.SendInstantMessage(event.UserPubkey, event.MessageText)
	if err != nil {
		c.onError(errors.New("failed to send instant message: " + err.Error()))
	}
}
