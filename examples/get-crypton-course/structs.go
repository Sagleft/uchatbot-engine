package main

import (
	"github.com/JulianToledano/goingecko"
	"github.com/Sagleft/uchatbot-engine"
)

type botApp struct {
	chatbot   *uchatbot.ChatBot
	coingecko *goingecko.Client
}
