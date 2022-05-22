package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/controllers"
)

type ControllerMain struct {
	controllers.Controller
	monitorHandler fiber.Handler
}

func NewControllerMain() ControllerMain {
	controller := ControllerMain{}

	return controller
}

func (сontroller ControllerMain) Monitor(config ...monitor.Config) fiber.Handler {
	handler := monitor.New(config...)

	сontroller.monitorHandler = func(connect *fiber.Ctx) error {
		request := controllers.NewRequest(connect).Admin()
		defer request.Store()

		if !request.Valid {
			return request.Err
		}

		return handler(connect)
	}

	return сontroller.monitorHandler
}

func (сontroller ControllerMain) Index(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect).Admin()
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	return connect.Status(200).Render("admin/main/index", fiber.Map{
		"Title":          "Admin Panel",
		"UrlAdminLogout": components.Route("admin.auth.logout", fiber.Map{}),
	}, "admin/layouts/main")
}

func (сontroller ControllerMain) NotFound(connect *fiber.Ctx) error {
	return connect.Status(404).Render("admin/errors/404", fiber.Map{
		"Title": "404",
		"Error": "Page Not Found. Error: 404",
	}, "admin/layouts/app")
}
