package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

type WebsocketClient struct {
	Connect   *websocket.Conn
	WsKey     uint64
	StackKey  uint64
	ClientKey uint64
}

func NewWebsocketClient(connect *websocket.Conn, wsKey uint64, stackKey uint64, clientKey uint64) (*WebsocketClient, error) {
	var err error

	if connect == nil {
		err = errors.New("No WebsocketStack")
	}

	wsClient := &WebsocketClient{
		Connect:   connect,
		WsKey:     wsKey,
		StackKey:  stackKey,
		ClientKey: clientKey,
	}

	return wsClient, err
}

func (wsClient *WebsocketClient) Key() string {
	return fmt.Sprintf("%d:%d:%d", wsClient.WsKey, wsClient.StackKey, wsClient.ClientKey)
}
