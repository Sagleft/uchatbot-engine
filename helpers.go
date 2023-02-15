package uchatbot

import (
	"errors"
	"reflect"

	"github.com/Sagleft/utopialib-go/v2/pkg/websocket"
)

type errorFunc func() error

func checkErrors(errChecks ...errorFunc) error {
	for _, errFunc := range errChecks {
		err := errFunc()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ChatBot) convertEventInterface(e interface{}, eventType string) (websocket.WsEvent, error) {
	// convert event interface
	event, isConvertable := e.(websocket.WsEvent)
	if !isConvertable {
		return websocket.WsEvent{}, errors.New("failed to convert " + eventType + " event interface. " +
			reflect.ValueOf(e).String() + " type received")
	}
	return event, nil
}

func ternaryInt(statement bool, a, b int) int {
	if statement {
		return a
	}
	return b
}
