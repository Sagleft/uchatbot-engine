# uchatbot-engine
Engine for creating chatbots for Utopia Messenger

## Concept

You don't want to understand Utopia API, but you have an idea how to make a bot that works with users in private and public messages.

The engine can:

1. process messages from contacts;
2. process messages in channels (private & public messages);
3. automatically logs into channels, can use the password for closed channels.

Planned features:
 - [ ] send welcome messages;
 - [ ] processing payments.

## Using the engine

1. Chatbots that raise and retain user activity in channels.
2. Bots for performing services to users.
3. Creating a bot constructor.

## Install

```bash
go get github.com/Sagleft/uchatbot-engine
```

## Example

```go
_, err := uchatbot.NewChatBot(uchatbot.ChatBotData{
    Client: &utopiago.UtopiaClient{
        Protocol: "http",
        Host:     APIHost,
        Token:    APIToken,
        Port:     22800,
        WsPort:   25000,
    },
    Chats: []uchatbot.Chat{
        {ID: "D53B4431FD604E2F0261792444797AA4"},
        {ID: "A59D8B62E1A59049564A4B0F8B457D45"},
    },
    Callbacks: uchatbot.ChatBotCallbacks{
        OnContactMessage:        OnContactMessage,
        OnChannelMessage:        OnChannelMessage,
        OnPrivateChannelMessage: OnPrivateChannelMessage,
    },
    UseErrorCallback: true,
    ErrorCallback:    onError,
})
if err != nil {
    log.Fatalln(err)
}
```

and set your callbacks. example:

```go
func OnContactMessage(m utopiago.InstantMessage) {
	fmt.Println("[CONTACT] " + m.Nick + ": " + m.Text)
}

func OnChannelMessage(m utopiago.WsChannelMessage) {
	fmt.Println("[CHANNEL] " + m.Nick + ": " + m.Text)
}

func OnPrivateChannelMessage(m utopiago.WsChannelMessage) {
	fmt.Println("[PRIVATE] " + m.Nick + ": " + m.Text)
}
```
