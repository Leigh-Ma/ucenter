package user

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
	"ucenter/controllers"
)

type loginController struct {
	controllers.ApiController
}

// GET /login/visitor
func (c *loginController) VisitorLogin() {
	f, resp := &form.FVisitorLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	user := models.GetUserByUuid(f.Uuid)
	if user.IsNew() {
		user.Insert(user)
	}

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, 24*3600)

	models.Upsert(auth)

	resp.Success(&http.D{"Token": auth.Token, "UserId": user.Id})

	c.RenderJson(resp)
}

// GET /login/token
func (c *loginController) TokenLogin() {
	f, resp := &form.FTokenLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	resp.Success(&http.D{
		"Token":  f.Token,
		"UserId": f.UserId,
	})

	player := models.GetPlayer(f.UserId)
	if player.IsNew() {
		resp.Error(http.ERR_USER_ID_INVALID)
		c.RenderJson(resp)
		return
	}

	auth := models.GetAuthToken(f.UserId)

	resp.Status(auth.VerifyToken(f.Token))

	c.RenderJson(resp)
}

// GET /login
func (c *loginController) Login() {
	f, resp := &form.FPasswordLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	user := models.GetUserByEmail(f.Email)
	if user.IsNew() {
		resp.Error(http.ERR_EMAIL_NOT_REGISTERED)
		c.RenderJson(resp)
		return
	}

	if !user.VerifyPassword(f.Password) {
		resp.Error(http.ERR_PASSWORD_ERROR)
		c.RenderJson(resp)
		return
	}

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, 24*3600)

	models.Upsert(auth)

	resp.Success(&http.D{"Token": auth.Token, "UserId": user.Id})

	c.RenderJson(resp)
}

func (c *loginController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"GET:  /visitor":    "VisitorLogin",
		"GET:  /token":      "TokenLogin",
		"GET:  /":           "Login",
	})
}
