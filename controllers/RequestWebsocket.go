package controllers

import (
	"strconv"
	"strings"
)

type RequestWebsocket struct {
	Ws       *Websocket
	WsStack  *WebsocketStack
	WsClient *WebsocketClient
	Message  string
}

func NewRequestWebsocket(ws *Websocket, wsStack *WebsocketStack, wsClient *WebsocketClient, message string) *RequestWebsocket {
	wsRequest := RequestWebsocket{
		Ws:       ws,
		WsStack:  wsStack,
		WsClient: wsClient,
		Message:  message,
	}

	return &wsRequest
}

func (wsRequest *RequestWebsocket) SendAll(message string) {
	wsRequest.Ws.SendAll(message)
}

func (wsRequest *RequestWebsocket) Send(key string, message string) {
	splitKey := strings.Split(key, ":")

	if len(splitKey) == 3 {
		var wsKey uint64
		var stackKey uint64
		var clientKey uint64

		if n, err := strconv.ParseUint(splitKey[0], 10, 64); err == nil {
			wsKey = n
			if n, err := strconv.ParseUint(splitKey[1], 10, 64); err == nil {
				stackKey = n
				if n, err := strconv.ParseUint(splitKey[2], 10, 64); err == nil {
					clientKey = n

					wsRequest.Ws.Send(wsKey, stackKey, clientKey, message)
				}
			}
		}
	}
}
