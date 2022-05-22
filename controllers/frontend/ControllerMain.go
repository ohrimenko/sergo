package frontend

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ohrimenko/sergo/controllers"
	"github.com/ohrimenko/sergo/models"
)

type ControllerMain struct {
	controllers.Controller
}

func NewControllerMain() ControllerMain {
	controller := ControllerMain{}

	return controller
}

func (сontroller ControllerMain) Hello(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	// Get value
	name := request.Sess.Get("name")

	// Set key/value
	request.Sess.Set("name", "john")

	user := models.NewUser()

	request.DB.Find(user, "id = ?", 12501)
	//components.Db.Raw("SELECT * FROM `users` WHERE id = ?", 25005).Scan(&user)

	if !user.Valid() {
		return connect.Status(404).Render("errors/404", fiber.Map{
			"Error":     "Error 404",
			"TextError": errors.New("User Not Found"),
		}, "layouts/main")
	}

	user.Gender.Set("male")
	user.Phone.Scan(nil)
	user.Phone.Set("0974721930")
	user.CoordinateLat.Set(50.450441)
	request.DB.Save(user)

	return connect.Status(200).Render("index", fiber.Map{
		"Title": name,
		"User":  user,
	}, "layouts/main")
}

func (сontroller ControllerMain) NotFound(connect *fiber.Ctx) error {
	return connect.Status(404).Render("errors/404", fiber.Map{
		"Title": "404",
		"Error": "Page Not Found. Error: 404",
	}, "layouts/main")
}
