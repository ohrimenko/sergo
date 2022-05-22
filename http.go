package main

import (
	"log"
	"time"

	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/config"
	"github.com/ohrimenko/sergo/models"
	"github.com/ohrimenko/sergo/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/utils"
	"github.com/joho/godotenv"

	"github.com/gofiber/template/html"

	"github.com/goccy/go-json"
)

var FiberApp *fiber.App

func main() {
	components.InitSerialize()

	// .env Variables validation
	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	_, err := components.DB()

	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")

	// Reload the templates on each render, good for development
	engine.Reload(false) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(false) // Optional. Default: false

	// Layout defines the variable name that is used to yield templates within layouts
	engine.Layout("embed") // Optional. Default: "embed"

	// Delims sets the action delimiters to the specified strings
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters

	engine.AddFunc("isAdmin", func(name models.User) bool {
		return true
	})

	// Fiber instance
	FiberApp = fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Sergo",
		AppName:       "Sergo App v1.0.1",
		Views:         engine,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	})

	components.FiberApp = FiberApp

	// Сжатие
	//FiberApp.Use(compress.New(compress.Config{
	//	Level: compress.LevelBestSpeed, // 1
	//}))

	FiberApp.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.Env("APP_KEY"),
	}))

	FiberApp.Use(favicon.New(favicon.Config{
		File: "./public/favicon.ico",
	}))

	FiberApp.Use(requestid.New(requestid.Config{
		Header:    "X-Custom-Header",
		Generator: utils.UUID,
	}))

	FiberApp.Static("/static", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	// Routes
	routes.App.Http(FiberApp)

	// Start server
	log.Fatal(FiberApp.Listen(":" + config.Env("HTTP_PORT")))
}
