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

	user, _ := models.UserM.GetVisitor(f.UUID)

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

	if err := models.TokenM.VerifyToken(f.UserId, f.Token); err != nil {
		resp.Error(err.Error())
		c.renderJson(resp)
		return
	}

	resp.Success(&http.D{"Token": f.Token, "UserId": f.UserId})

	c.renderJson(resp)
}

// GET /login
func (c *LoginController) Login() {
	f, resp := &form.FPasswordLogin{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	user, isNew := models.UserM.GetByEmail(f.Email, f.UUID)
	if isNew {
		resp.PasswordError("User not found")
		c.renderJson(resp)
		return
	}

	if !user.VerifyPassword(f.Password) {
		resp.PasswordError()
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
		"GET:  /visitor":     "VisitorLogin",
		"GET:  /tokenLogin":  "TokenLogin",
		"GET:  /":            "Login",
	})
}
