package api

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type LoginController struct {
	apiController
}

// GET /login/visitor
func (c *LoginController) VisitorLogin() {
	f, resp := &form.FVisitorLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
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

	c.renderJson(resp)
}

// GET /login/token
func (c *LoginController) TokenLogin() {
	f, resp := &form.FTokenLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	resp.Success(&http.D{
		"Token":  f.Token,
		"UserId": f.UserId,
	})

	player := models.GetPlayer(f.UserId)
	if player.IsNew() {
		resp.Error(http.ERR_USER_ID_INVALID)
		c.renderJson(resp)
		return
	}

	auth := models.GetAuthToken(f.UserId)

	resp.Set(auth.VerifyToken(f.Token))

	c.renderJson(resp)
}

// GET /login
func (c *LoginController) Login() {
	f, resp := &form.FPasswordLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	user := models.GetUserByEmail(f.Email)
	if user.IsNew() {
		resp.Error(http.ERR_EMAIL_NOT_REGISTERED)
		c.renderJson(resp)
		return
	}

	if !user.VerifyPassword(f.Password) {
		resp.Error(http.ERR_PASSWORD_ERROR)
		c.renderJson(resp)
		return
	}

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, 24*3600)

	models.Upsert(auth)

	resp.Success(&http.D{"Token": auth.Token, "UserId": user.Id})

	c.renderJson(resp)
}

func (c *LoginController) Export() func(string) {
	return export(c, map[string]string{
		"GET:  /visitor":    "VisitorLogin",
		"GET:  /tokenLogin": "TokenLogin",
		"GET:  /":           "Login",
	})
}
