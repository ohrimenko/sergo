package routes

import (
	"log"
	"time"

	"github.com/ohrimenko/sergo/controllers"
	"github.com/ohrimenko/sergo/controllers/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (route Route) Websocket(app *fiber.App) {
	wsControllerMain := ws.NewControllerMain()
	ws := controllers.NewWebsocket(1000000, wsControllerMain.OnConnect, wsControllerMain.OnMessage, wsControllerMain.OnClose)

	app.Static("/", "./public/websocket/home.html", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
		}()

		wsClient, err := ws.NewWebsocketClient(c)

		if err == nil {
			// When the function returns, unregister the client and close the connection
			defer func() {
				ws.Unregister(wsClient)
			}()

			// Register the client
			ws.Register(wsClient)

			for {
				messageType, message, err := c.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Println("read error:", err)
					}

					return // Calls the deferred function, i.e. closes the connection on error
				}

				if messageType == websocket.TextMessage {
					// Broadcast the received message
					ws.Broadcast(wsClient, string(message))
				} else {
					log.Println("websocket message received of type", messageType)
				}
			}
		}
	})).Name("main.websocket.index")

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	}).Name("main.control.socket")
}
