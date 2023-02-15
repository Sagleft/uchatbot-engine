package uchatbot

import (
	"errors"
	"reflect"
)

type sendMessageTask struct {
	UserPubkey  string
	MessageText string
}

type sendChannelPrivateMessageTask struct {
	ChannelID      string
	UserPubkeyHash string
	MessageText    string
}

// SetReadonly - enable or disable channel readonly mode
func (c *ChatBot) SetReadonly(channelID string, readOnly bool) error {
	return c.client.EnableReadOnly(channelID, readOnly)
}

// SendContactMessage - send message to contact.
// it works with queue (buffer)
func (c *ChatBot) SendContactMessage(userPubkey string, msgText string) {
	if c.rateLimiters.InstantMessage.Enabled {
		c.rateLimiters.InstantMessage.L.Wait()
	}

	c.queues.InstantMessages.AddEvent(sendMessageTask{
		UserPubkey:  userPubkey,
		MessageText: msgText,
	})
}

func (c *ChatBot) handleSendInstantMessageTask(e interface{}) {
	event, isConvertable := e.(sendMessageTask)
	if !isConvertable {
		c.onError(errors.New("failed to convert send message task: " +
			reflect.ValueOf(e).String() + " type received"))
		return
	}

	_, err := c.client.SendInstantMessage(event.UserPubkey, event.MessageText)
	if err != nil {
		c.onError(errors.New("failed to send instant message: " + err.Error()))
	}
}

// SendChannelPrivateMessage - send message to contact in channel (in private chat).
// it works with queue (buffer)
func (c *ChatBot) SendChannelPrivateMessage(channel, userPubkeyHash, msgText string) {
	if c.rateLimiters.ChannelPrivateMessage.Enabled {
		c.rateLimiters.ChannelPrivateMessage.L.Wait()
	}

	c.queues.SendPrivateChannelMessage.AddEvent(sendChannelPrivateMessageTask{
		ChannelID:      channel,
		UserPubkeyHash: userPubkeyHash,
		MessageText:    msgText,
	})
}

func (c *ChatBot) handleSendPrivateChannelMessageTask(e interface{}) {
	event, isConvertable := e.(sendChannelPrivateMessageTask)
	if !isConvertable {
		c.onError(errors.New("failed to convert send private channel message task: " +
			reflect.ValueOf(e).String() + " type received"))
		return
	}

	_, err := c.client.SendChannelContactMessage(event.ChannelID, event.UserPubkeyHash, event.MessageText)
	if err != nil {
		c.onError(errors.New("failed to send channel private message: " + err.Error()))
	}
}

// GetOwnPubkey - get account public key
func (c *ChatBot) GetOwnPubkey() (string, error) {
	data, err := c.client.GetOwnContact()
	if err != nil {
		return "", err
	}

	return data.Pubkey, nil
}
