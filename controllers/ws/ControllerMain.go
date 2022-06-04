package ws

import (
	"log"

	"github.com/ohrimenko/sergo/controllers"
)

type ControllerMain struct {
	controllers.Controller
}

func NewControllerMain() ControllerMain {
	controller := ControllerMain{}

	return controller
}

func (сontroller ControllerMain) Register(wsClient *controllers.WebsocketClient) {
	log.Println("connection registered: ", wsClient.Key())
}

func (сontroller ControllerMain) Message(wsRequest *controllers.RequestWebsocket) {
	wsRequest.Send(wsRequest.WsClient.Key(), "send...")
	wsRequest.SendAll(wsRequest.Message)
	log.Println("Send: ", wsRequest.WsClient.Key(), " - ", wsRequest.Message)
}

func (сontroller ControllerMain) Unregister(wsClient *controllers.WebsocketClient) {
	log.Println("connection unregistered: ", wsClient.Key())
}
