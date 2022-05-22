package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohrimenko/sergo/controllers/admin"
	"github.com/ohrimenko/sergo/controllers/frontend"
)

func (route Route) Http(app *fiber.App) {
	app.Get("/admin", admin.NewControllerMain().Index).Name("admin.index")
	app.Get("/admin/login", admin.NewControllerAuth().Login).Name("admin.auth.login")
	app.Post("/admin/authorize", admin.NewControllerAuth().Authorize).Name("admin.auth.authorize")
	app.Get("/admin/logout", admin.NewControllerAuth().Logout).Name("admin.auth.logout")
	app.Get("/admin/monitor", admin.NewControllerMain().Monitor()).Name("admin.auth.monitor")
	// 404 Handler Admin
	app.Use("/admin*", admin.NewControllerMain().NotFound).Name("admin.404")

	app.Get("/", frontend.NewControllerMain().Hello).Name("main.index")
	app.Get("/login", frontend.NewControllerAuth().Login).Name("main.auth.login")
	app.Get("/authorize", frontend.NewControllerAuth().Authorize).Name("main.auth.authorize")
	app.Get("/logout", frontend.NewControllerAuth().Logout).Name("main.auth.logout")

	// 404 Handler Main
	app.Use(frontend.NewControllerMain().NotFound).Name("main.404")
}
