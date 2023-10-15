package main

import (
	"github.com/Sagleft/uchatbot-engine"
)

// just for example.
// move these constants to the configuration file
// or environment variables in the production application

const APIToken = "your-utopia-api-token"
const utopiaHost = "127.0.0.1"
const utopiaPort = 20000
const utopiaWsPort = 25000
const botCommand = "show crypton rate"
const defaultErrorMessage = "I'm having problems and can't process the request"

type destinationType int

const (
	destTypeChannel destinationType = iota
	destTypeContact
)

var chats = []uchatbot.Chat{
	{ID: "D53B4431FD604E2F0261792444797AA4"},
}
