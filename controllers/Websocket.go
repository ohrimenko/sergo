package controllers

import (
	"log"
	"math/rand"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type ClientsStack struct {
	Clients map[uint64]*WebsocketClient
}

type WebsocketStack struct {
	Mutex      sync.Mutex
	Ws         *Websocket
	Register   chan *WebsocketClient
	Broadcast  chan *BroadcastWebsocket
	Unregister chan *WebsocketClient
	Count      uint64
	I          uint64
	Key        uint64
}

type Websocket struct {
	Mutex           sync.Mutex
	Stack           map[uint64]*WebsocketStack
	Key             uint64
	MaxCountInStack uint64
	I               uint64
	FuncRegister    func(wsClient *WebsocketClient)
	FuncMessage     func(wsRequest *WebsocketClient, message string)
	FuncUnregister  func(wsClient *WebsocketClient)
}

var Websockets map[uint64]*Websocket

func (stack *ClientsStack) DeleteClient(ClientKey uint64) {
	if _, ok := stack.Clients[ClientKey]; ok {
		delete(stack.Clients, ClientKey)
	}
}

func (ws *Websocket) DeleteStack(key uint64) {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	if _, ok := ws.Stack[key]; ok {
		delete(ws.Stack, key)
	}
}

func (wsStack *WebsocketStack) CountIncrement() {
	wsStack.Mutex.Lock()
	defer wsStack.Mutex.Unlock()
	wsStack.Count++
}

func (wsStack *WebsocketStack) CountDecrement() {
	wsStack.Mutex.Lock()
	defer wsStack.Mutex.Unlock()
	if wsStack.Count > 1 {
		wsStack.Count--
	} else {
		wsStack.Count = 0
	}
}

func (wsStack *WebsocketStack) Delete() {
	if wsStack.Ws != nil {
		wsStack.Ws.DeleteStack(wsStack.Key)
	}
}

func NewWebsocket(maxConnect uint64, register func(wsClient *WebsocketClient), message func(wsRequest *WebsocketClient, message string), unregister func(wsClient *WebsocketClient)) *Websocket {
	if Websockets == nil {
		Websockets = make(map[uint64]*Websocket)
	}

	key := uint64(rand.Intn(1000000) + 1000000)

	Websockets[key] = &Websocket{
		Stack:           map[uint64]*WebsocketStack{},
		Key:             key,
		MaxCountInStack: 5000,
		I:               0,
		FuncRegister:    register,
		FuncMessage:     message,
		FuncUnregister:  unregister,
	}

	countStack := int(maxConnect / 5000)

	if countStack < 50 {
		countStack = 50
	}

	for i := 0; i < countStack; i++ {
		Websockets[key].NewWebsocketStack()
	}

	return Websockets[key]
}

func (wsStack *WebsocketStack) RunHub() {
	stack := ClientsStack{
		Clients: make(map[uint64]*WebsocketClient),
	}

	for {
		select {
		case wsClient := <-wsStack.Register:
			wsStack.CountIncrement()
			stack.Clients[wsClient.ClientKey] = wsClient

		case wsBroadcast := <-wsStack.Broadcast:
			// Send the message to all clients
			if wsBroadcast.ClientKey == 0 {
				for _, wsClient := range stack.Clients {
					if err := wsClient.Connect.WriteMessage(websocket.TextMessage, []byte(wsBroadcast.Message)); err != nil {
						log.Println("write error:", err)

						wsClient.Connect.WriteMessage(websocket.CloseMessage, []byte{})
						wsClient.Connect.Close()

						stack.DeleteClient(wsClient.ClientKey)

						wsStack.CountDecrement()
					}
				}
			} else if wsClient, ok := stack.Clients[wsBroadcast.ClientKey]; ok {
				if err := wsClient.Connect.WriteMessage(websocket.TextMessage, []byte(wsBroadcast.Message)); err != nil {
					log.Println("write error:", err)

					wsClient.Connect.WriteMessage(websocket.CloseMessage, []byte{})
					wsClient.Connect.Close()

					stack.DeleteClient(wsClient.ClientKey)

					wsStack.CountDecrement()
				}
			}

		case wsClient := <-wsStack.Unregister:
			// Remove the client from the hub
			stack.DeleteClient(wsClient.ClientKey)

			wsStack.CountDecrement()
		}
	}
}

func (ws *Websocket) NewWebsocketClient(connection *websocket.Conn) (*WebsocketClient, error) {
	var stack *WebsocketStack
	var wsClient *WebsocketClient
	var err error

	if _, ok := ws.Stack[1]; ok {
		stack = ws.Stack[1]

		for _, wsStack := range ws.Stack {
			if wsStack.Count < stack.Count {
				stack = wsStack
			}
		}
	}

	wsClient, err = stack.NewWebsocketClient(connection)

	return wsClient, err
}

func (ws *Websocket) NewWebsocketStack() *WebsocketStack {
	ws.I++
	ws.Stack[ws.I] = &WebsocketStack{
		Ws:         ws,
		Register:   make(chan *WebsocketClient),
		Broadcast:  make(chan *BroadcastWebsocket),
		Unregister: make(chan *WebsocketClient),
		Count:      0,
		I:          0,
		Key:        ws.I,
	}

	go ws.Stack[ws.I].RunHub()

	return ws.Stack[ws.I]
}

func (wsStack *WebsocketStack) NewWebsocketClient(connection *websocket.Conn) (*WebsocketClient, error) {
	var err error
	var wsClient *WebsocketClient

	wsStack.Mutex.Lock()
	wsStack.I++

	I := wsStack.I
	wsStack.Mutex.Unlock()

	wsClient, err = NewWebsocketClient(connection, wsStack.Ws.Key, wsStack.Key, I)

	return wsClient, err
}

func (ws *Websocket) Register(wsClient *WebsocketClient) {
	if wsStack, ok := ws.Stack[wsClient.StackKey]; ok {
		wsStack.Register <- wsClient

		ws.FuncRegister(wsClient)
	}
}

func (ws *Websocket) Broadcast(wsClient *WebsocketClient, s string) {
	if _, ok := ws.Stack[wsClient.StackKey]; ok {
		ws.FuncMessage(wsClient, s)
	}
}

func (ws *Websocket) Unregister(wsClient *WebsocketClient) {
	if wsStack, ok := ws.Stack[wsClient.StackKey]; ok {
		wsStack.Unregister <- wsClient

		ws.FuncUnregister(wsClient)
	}
}

func (ws *Websocket) SendAll(message string) {
	wsBroadcast := NewBroadcastWebsocket(0, 0, 0, message)

	for _, wsStack := range ws.Stack {
		if wsStack.Count > 0 {
			wsStack.Broadcast <- wsBroadcast
		}
	}
}

func (ws *Websocket) Send(wsKey uint64, stackKey uint64, clientKey uint64, message string) {
	if wsKey > 0 && stackKey > 0 && clientKey > 0 && ws.Key == wsKey {
		if wsStack, ok := ws.Stack[stackKey]; ok {
			if wsStack.Count > 0 {
				wsBroadcast := NewBroadcastWebsocket(wsKey, stackKey, clientKey, message)

				wsStack.Broadcast <- wsBroadcast
			}
		}
	}
}
