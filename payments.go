package uchatbot

import (
	"github.com/Sagleft/utopialib-go/v2/pkg/structs"
)

const (
	CurrencyCrypton CurrencyType = "CRP"
	CurrencyUUSD    CurrencyType = "UUSD"
)

type CurrencyType string

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
