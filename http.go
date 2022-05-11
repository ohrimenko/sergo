package main

import (
	"log"

	"github.com/ohrimenko/sergo/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Fiber instance
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})

	// Routes
	app.Get("/", hello)

	// .env Variables validation
	if err := godotenv.Load("./.env"); err != nil {
		panic("Error loading .env file")
	}

	// Start server
	log.Fatal(app.Listen(config.Env("HTTP_PORT")))
}

// Handler
func hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
