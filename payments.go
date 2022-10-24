package uchatbot

import (
	utopiago "github.com/Sagleft/utopialib-go"
)

const (
	CurrencyCrypton CurrencyType = "CRP"
	CurrencyUUSD    CurrencyType = "UUSD"
)

type CurrencyType string

// SendCoins from main account
func (c *ChatBot) SendCoins(currency CurrencyType, pubkey string, amount float64) {
	c.data.Client.SendPayment(utopiago.SendPaymentTask{
		To:          pubkey,
		Amount:      amount,
		CurrencyTag: string(currency),
	})
}

// SendCoins from crypto card
func (c *ChatBot) SendCoinsFromCard(currency CurrencyType, pubkey string, amount float64, fromCard string) {
	c.data.Client.SendPayment(utopiago.SendPaymentTask{
		To:          pubkey,
		Amount:      amount,
		CurrencyTag: string(currency),
		FromCardID:  fromCard,
	})
}
