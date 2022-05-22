package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
}

func (сontroller Controller) Abort500(connect *fiber.Ctx) error {
	if connect != nil {
		return connect.Status(500).Render("errors/500", fiber.Map{
			"Error":     "Error 500",
			"TextError": "500 Internal Server Error",
		}, "layouts/main")
	} else {
		return errors.New("Connect is <nil>")
	}
}

func (сontroller Controller) Abort404(connect *fiber.Ctx) error {
	if connect != nil {
		return connect.Status(404).Render("errors/500", fiber.Map{
			"Error":     "Error 500",
			"TextError": "500 Internal Server Error",
		}, "layouts/main")
	} else {
		return errors.New("Connect is <nil>")
	}
}
