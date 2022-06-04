package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/ohrimenko/sergo/controllers/admin"
	"github.com/ohrimenko/sergo/controllers/frontend"
)

func (route Route) Http(app *fiber.App) {
	adminControllerMain := admin.NewControllerMain()
	adminControllerAuth := admin.NewControllerAuth()

	frontendControllerMain := frontend.NewControllerMain()
	frontendControllerAuth := frontend.NewControllerAuth()

	app.Use(favicon.New(favicon.Config{
		File: "./public/favicon.ico",
	}))

	app.Static("/static", "./public", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	app.Get("/admin", adminControllerMain.Index).Name("admin.index")
	app.Get("/admin/login", adminControllerAuth.Login).Name("admin.auth.login")
	app.Post("/admin/authorize", adminControllerAuth.Authorize).Name("admin.auth.authorize")
	app.Get("/admin/logout", adminControllerAuth.Logout).Name("admin.auth.logout")
	app.Get("/admin/monitor", adminControllerMain.Monitor()).Name("admin.auth.monitor")
	// 404 Handler Admin
	app.Use("/admin*", adminControllerMain.NotFound).Name("admin.404")

	app.Get("/", frontendControllerMain.Hello).Name("main.index")
	app.Get("/login", frontendControllerAuth.Login).Name("main.auth.login")
	app.Get("/authorize", frontendControllerAuth.Authorize).Name("main.auth.authorize")
	app.Get("/logout", frontendControllerAuth.Logout).Name("main.auth.logout")

	// 404 Handler Main
	app.Use(frontendControllerMain.NotFound).Name("main.404")
}
