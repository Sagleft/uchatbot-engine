package main

import (
	"errors"
	"fmt"

	"github.com/JulianToledano/goingecko"
)

func crateCoingeckoClient() *goingecko.Client {
	return goingecko.NewClient(nil)
}

func getCryptonRate(client *goingecko.Client) (float64, error) {
	priceData, err := client.SimplePrice("crypton", "usd", false, false, false, false)
	if err != nil {
		return 0, fmt.Errorf("get Crypton to USD price: %w", err)
	}
	if priceVal, isExists := priceData["crypton"]; isExists {
		if price, isExists := priceVal["usd"]; isExists {
			return price, nil
		}
		return 0, errors.New("crypton rate not found")
	} else {
		return 0, errors.New("crypton rate not found")
	}
}
