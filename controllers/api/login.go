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

	user, _ := models.UserM.GetVisitor(f.Uuid)

	auth := models.TokenM.GetUserToken(user.Id)

	auth.SetNewToken(user.Id, 24*3600)

	resp.Success(&http.D{"Token": auth.Token, "UserId": user.Id})

	c.renderJson(resp)
}

// GET /login/token
func (c *LoginController) TokenLogin() {
	f, resp := &form.FTokenLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	status := models.TokenM.VerifyToken(f.UserId, f.Token)

	resp.Set(status, &http.D{
		"Token":  f.Token,
		"UserId": f.UserId,
	})

	c.renderJson(resp)
}

// GET /login
func (c *LoginController) Login() {
	f, resp := &form.FPasswordLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	user, isNew := models.UserM.GetByEmail(f.Email, f.Uuid)
	if isNew {
		resp.Error(http.ERR_EMAIL_NOT_REGISTERED)
		c.renderJson(resp)
		return
	}

	if !user.VerifyPassword(f.Password) {
		resp.Error(http.ERR_PASSWORD_ERROR)
		c.renderJson(resp)
		return
	}

	auth := models.TokenM.GetUserToken(user.Id)

	auth.SetNewToken(user.Id, 24*3600)

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
