package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohrimenko/sergo/controllers"
)

type ControllerMain struct {
	controllers.Controller
}

func NewControllerMain() ControllerMain {
	controller := ControllerMain{}

	return controller
}
