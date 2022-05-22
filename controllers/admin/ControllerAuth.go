package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/controllers"
	"github.com/ohrimenko/sergo/models"
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

	return connect.Status(200).Render("admin/auth/login", fiber.Map{
		"Title":         "Authorize",
		"UrlAdminLogin": components.Route("admin.auth.authorize", fiber.Map{}),
		"Request":       request,
	}, "admin/layouts/app")
}

func (сontroller ControllerAuth) Authorize(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	type AuthorizeInput struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var input AuthorizeInput

	if err := connect.BodyParser(&input); err != nil {
		return сontroller.Abort404(connect)
	}

	user := models.NewUser()

	request.DB.Where("`login` = ? ", input.Login).Find(user)

	if user.Valid() {
		if components.CheckPasswordHash(input.Password, user.Password.Get()) {
			request.Sess.Set("AuthUserId", user.Id.Get())

			return connect.Status(302).RedirectToRoute("admin.index", fiber.Map{})
		} else {
			request.Error("Password", "Не верный пароль")
		}
	} else {
		request.Error("Login", "Пользователь не найден")
	}

	request.Old("Login", input.Login)
	request.Old("Password", input.Password)

	return connect.Status(302).RedirectToRoute("admin.auth.login", fiber.Map{})
}

func (сontroller ControllerAuth) Logout(connect *fiber.Ctx) error {
	request := controllers.NewRequest(connect)
	defer request.Store()

	if !request.Valid {
		return request.Err
	}

	request.Sess.Delete("AuthUserId")

	return connect.Status(302).RedirectToRoute("admin.auth.login", fiber.Map{})
}
