package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/ohrimenko/sergo/components"
	"github.com/ohrimenko/sergo/models"
	"gorm.io/gorm"
)

type Request struct {
	Connect      *fiber.Ctx
	DB           *gorm.DB
	Sess         *session.Session
	User         *models.User
	Err          error
	Valid        bool
	ValidDB      bool
	ValidSession bool
	Messages     components.MessagesMap
	Errors       components.MessagesMap
	Olds         components.MessagesMap
}

func NewRequest(connect *fiber.Ctx) *Request {
	req := Request{}

	req.Connect = connect

	req.Valid = true
	req.ValidDB = false
	req.ValidSession = false

	req.DB, req.Err = components.DB()

	if req.Err != nil {
		req.Valid = false
	} else {
		req.ValidDB = true
	}

	req.User = models.NewUser()

	if req.Connect == nil {
		req.Valid = false
		req.Err = errors.New("User Not Found")
	} else {
		req.Sess, req.Err = components.Session().Get(req.Connect)
		if req.Err == nil {
			req.ValidSession = true
			if req.Valid {
				switch value := req.Sess.Get("AuthUserId").(type) {
				case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
					req.DB.Find(req.User, "id = ?", value)
				}
			}

			if value, ok := req.Sess.Get("Errors").(components.MessagesMap); ok {
				req.Errors = value
			}

			if value, ok := req.Sess.Get("Messages").(components.MessagesMap); ok {
				req.Messages = value
			}

			if value, ok := req.Sess.Get("Olds").(components.MessagesMap); ok {
				req.Olds = value
			}
		} else {
			req.Valid = false
		}
	}

	if req.Messages == nil {
		req.Messages = components.MessagesMap{}
	}

	if req.Errors == nil {
		req.Errors = components.MessagesMap{}
	}

	if req.Olds == nil {
		req.Olds = components.MessagesMap{}
	}

	if !req.Valid {
		if req.Connect != nil {
			req.Err = req.Connect.Status(500).Render("errors/500", fiber.Map{
				"Error":     "Error 500",
				"TextError": req.Err,
			}, "layouts/main")
		}
	}

	return &req
}

func (req *Request) IsAuth() bool {
	if req.Valid {
		if req.User != nil {
			if req.User.Valid() {
				return true
			}
		}
	}

	return false
}

func (req *Request) Auth() *Request {
	if req.Valid {
		if !req.IsAuth() {
			req.Valid = false
		}

		if !req.Valid {
			if req.Connect != nil {
				req.Err = req.Connect.Status(302).RedirectToRoute("main.auth.login", fiber.Map{})
			}
		}
	}

	return req
}

func (req *Request) IsAdmin() bool {
	if req.Valid {
		if req.User != nil {
			if req.User.Valid() {
				if req.User.Type.Get() == 1 {
					return true
				}
			}
		}
	}

	return false
}

func (req *Request) Admin() *Request {
	if req.Valid {
		if !req.IsAdmin() {
			req.Valid = false
		}

		if !req.Valid {
			if req.Connect != nil {
				req.Err = req.Connect.Status(302).RedirectToRoute("admin.auth.login", fiber.Map{})
			}
		}
	}

	return req
}

func (req *Request) Error(key string, value ...any) string {
	if req.Valid {
		if len(value) > 0 {
			req.Errors[key] = &components.Message{}

			for _, v := range value {
				req.Errors[key].Scan(v)
			}
		} else if val, ok := req.Errors[key]; ok {
			return val.Get()
		}
	}

	return ""
}

func (req *Request) Message(key string, value ...any) string {
	if req.Valid {
		if len(value) > 0 {
			req.Errors[key] = &components.Message{}

			for _, v := range value {
				req.Errors[key].Scan(v)
			}
		} else if val, ok := req.Errors[key]; ok {
			return val.Get()
		}
	}

	return ""
}

func (req *Request) Old(key string, value ...any) string {
	if req.Valid {
		if len(value) > 0 {
			req.Olds[key] = &components.Message{}

			for _, v := range value {
				req.Olds[key].Scan(v)
			}
		} else if val, ok := req.Olds[key]; ok {
			return val.Get()
		}
	}

	return ""
}

func (req *Request) OldOrValue(key string, value ...any) string {
	res := ""

	if req.Valid {
		if val, ok := req.Olds[key]; ok {
			res = val.String()
		}

		if res == "" && len(value) > 0 {
			for _, v := range value {
				components.СonvertAssign(&res, v)

				if res != "" {
					break
				}
			}
		}
	}

	return res
}

func (req *Request) Store() *Request {
	if req.Valid {
		if req.ValidSession {
			if req.Messages != nil {
				for key, message := range req.Messages {
					if !message.IsSave() {
						delete(req.Messages, key)
					}
				}

				if len(req.Messages) > 0 {
					req.Sess.Set("Messages", req.Messages)
				} else {
					req.Sess.Delete("Messages")
				}
			} else {
				req.Sess.Delete("Messages")
			}
			if req.Errors != nil {
				for key, message := range req.Errors {
					if !message.IsSave() {
						delete(req.Errors, key)
					}
				}

				if len(req.Errors) > 0 {
					req.Sess.Set("Errors", req.Errors)
				} else {
					req.Sess.Delete("Errors")
				}
			} else {
				req.Sess.Delete("Errors")
			}
			if req.Olds != nil {
				for key, message := range req.Olds {
					if !message.IsSave() {
						delete(req.Olds, key)
					}
				}

				if len(req.Olds) > 0 {
					req.Sess.Set("Olds", req.Olds)
				} else {
					req.Sess.Delete("Olds")
				}
			} else {
				req.Sess.Delete("Olds")
			}

			req.Sess.Save()
		}
	}

	return req
}
