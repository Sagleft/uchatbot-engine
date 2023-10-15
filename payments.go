package uchatbot

import (
	"errors"
	"fmt"

	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

const (
	CurrencyCrypton CurrencyType = "CRP"
	CurrencyUUSD    CurrencyType = "UUSD"
)

const defaultDonateUCodeFilename = "channel_owner_pubkey.jpg"
const defaultDonateMessage = "Hi all! Support the author of this channel! " +
	"Donate coins to the author:"

type CurrencyType string

type DonateService struct {
	c *ChatBot

	channelID    string
	message      string
	uCodeEnabled bool
	uCodeComment string
}

// SendCoins from main account
func (c *ChatBot) SendCoins(currency CurrencyType, pubkey string, amount float64) {
	c.client.SendPayment(structs.SendPaymentTask{
		To:          pubkey,
		Amount:      amount,
		CurrencyTag: string(currency),
	})
}

// SendCoins from crypto card
func (c *ChatBot) SendCoinsFromCard(currency CurrencyType, pubkey string, amount float64, fromCard string) {
	c.client.SendPayment(structs.SendPaymentTask{
		To:          pubkey,
		Amount:      amount,
		CurrencyTag: string(currency),
		FromCardID:  fromCard,
	})
}

// RequestDonate - ask chat users to support the channel author
func (c *ChatBot) RequestDonate(channelID string) *DonateService {
	return &DonateService{
		c:         c,
		channelID: channelID,
		message:   defaultDonateMessage,
	}
}

// EnableUCode - add uCode with address to message
func (srv *DonateService) EnableUCode(enabled bool, comment string) *DonateService {
	srv.uCodeEnabled = enabled
	srv.uCodeComment = comment
	return srv
}

// SetMessage - set a custom message that will appear before the payment address
func (srv *DonateService) SetMessage(newMessage string) *DonateService {
	srv.message = newMessage
	return srv
}

// GetDonateMessage returns ownerPubkey, message, error
func (srv *DonateService) GetDonateMessage() (string, string, error) {
	if srv.channelID == "" {
		return "", "", errors.New("channel ID is not set")
	}

	// find channel author payment address
	channelData, err := srv.c.GetClient().GetChannelInfo(srv.channelID)
	if err != nil {
		return "", "", fmt.Errorf(
			"get channel %q data: %w",
			srv.channelID, err,
		)
	}

	ownerPubkey := channelData.Owner
	if ownerPubkey == "" {
		return "", "", fmt.Errorf(
			"channel %q owner unknown: pubkey is empty",
			srv.channelID,
		)
	}

	msg := fmt.Sprintf(
		"%s %s",
		srv.message, ownerPubkey,
	)
	return ownerPubkey, msg, nil
}

func (srv *DonateService) Do() error {
	ownerPubkey, msg, err := srv.GetDonateMessage()
	if err != nil {
		return fmt.Errorf("get donate message: %w", err)
	}

	if err := srv.c.SendChannelMessage(srv.channelID, msg); err != nil {
		return fmt.Errorf(
			"send donate message to channel %q: %w",
			srv.channelID, err,
		)
	}

	if srv.uCodeEnabled {
		ucodeBytes, err := srv.c.client.UCodeEncode(ownerPubkey, "BASE64", "JPG", 128)
		if err != nil {
			return fmt.Errorf("encode channel owner pubkey to uCode: %w", err)
		}

		if _, err := srv.c.client.SendChannelPicture(
			srv.channelID,
			ucodeBytes,
			srv.uCodeComment,
			defaultDonateUCodeFilename,
		); err != nil {
			return fmt.Errorf(
				"send ucode to channel %q: %w",
				srv.channelID, err,
			)
		}
	}
	return nil
}
