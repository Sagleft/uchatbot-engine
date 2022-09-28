package uchatbot

// SetReadonly - enable or disable channel readonly mode
func (c *ChatBot) SetReadonly(channelID string, readOnly bool) error {
	return c.data.Client.EnableReadOnly(channelID, readOnly)
}
