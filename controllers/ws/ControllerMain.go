package ws

import (
	"fmt"

	"github.com/ohrimenko/sergo/controllers"
)

type ControllerMain struct {
	controllers.Controller
}

func NewControllerMain() ControllerMain {
	controller := ControllerMain{}

	return controller
}

func (сontroller ControllerMain) OnConnect(wsClient *controllers.WebsocketClient) {
	wsClient.SendAll(fmt.Sprint("connection registered: ", wsClient.Key()))
	wsClient.Key()
}

func (сontroller ControllerMain) OnMessage(wsClient *controllers.WebsocketClient, message string) {
	wsClient.Send(wsClient.Key(), "send...")
	wsClient.SendAll(message)
	wsClient.Key()
}

func (сontroller ControllerMain) OnClose(wsClient *controllers.WebsocketClient) {
	wsClient.SendAll(fmt.Sprint("connection unregistered: ", wsClient.Key()))
	wsClient.Key()
}
