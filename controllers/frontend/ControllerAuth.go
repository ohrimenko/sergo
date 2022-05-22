package frontend

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/controllers"
)

type ControllerAuth struct {
	controllers.Controller
}

func NewControllerAuth() ControllerAuth {
	controller := ControllerAuth{}

	return controller
}

func (сontroller ControllerAuth) Login(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	return connect.Status(200).Render("auth/login", fiber.Map{
		"Title":            "Authorize",
		"UrlFrontendLogin": components.Route("main.auth.authorize", fiber.Map{}),
	}, "layouts/main")
}

func (сontroller ControllerAuth) Authorize(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	return connect.Status(302).RedirectToRoute("main.auth.login", fiber.Map{})
}

func (сontroller ControllerAuth) Logout(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	request.Sess.Delete("AuthUserId")

	return connect.Status(302).RedirectToRoute("main.auth.login", fiber.Map{})
}
