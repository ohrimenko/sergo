package main

import (
	"log"
	"time"

	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/config"
	"github.com/ohrimenko/sergo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/utils"
	"github.com/joho/godotenv"

	"github.com/gofiber/template/html"

	"github.com/goccy/go-json"

	"math/rand"
)

var FiberApp *fiber.App

func main() {
	rand.Seed(time.Now().UnixNano())

	components.InitSerialize()

	// .env Variables validation
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err := components.DB()

	if err != nil {
		log.Fatal(err)
	}

	defer components.CloseDB()

	engine := html.New("./views", ".html")

	// Reload the templates on each render, good for development
	engine.Reload(false) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(false) // Optional. Default: false

	// Layout defines the variable name that is used to yield templates within layouts
	engine.Layout("embed") // Optional. Default: "embed"

	// Delims sets the action delimiters to the specified strings
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	// Fiber instance
	FiberApp = fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Server-Http",
		AppName:       "Server-Http App v1.0.1",
		Views:         engine,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	components.FiberApp = FiberApp

	FiberApp.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.Env("APP_KEY"),
	}))

	FiberApp.Use(requestid.New(requestid.Config{
		Header:    "X-Custom-Header",
		Generator: utils.UUID,
	}))

	// Routes
	routes.App.Http(FiberApp)

	// Start server
	log.Fatal(FiberApp.Listen(":" + config.Env("HTTP_PORT")))
}
