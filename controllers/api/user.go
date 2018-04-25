package api

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type UserController struct {
	apiController
}

// POST:  /user/modify_pwd
func (c *UserController) ModifyPassword() {
	f, resp := &form.FModifyPassword{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	if f.Password == f.PasswordNew {
		c.renderJson(resp.Error(http.ERR_PASSWORD_NOT_CHANGED))
		return
	}

	user := models.GetUserByEmail(f.Email)
	if user.IsNew() {
		c.renderJson(resp.Error(http.ERR_EMAIL_NOT_REGISTERED))
		return
	}

	user.SetPassword(f.Password)
	if _, err := user.Update(user); err != nil {
		c.renderJson(resp.Error(http.ERR_DATA_BASE_ERROR))
		return
	}

	models.Upsert(user)

	//expire token
	c.renderJson(resp.Success(&http.D{
		"UserId": user.Id,
		"Email": user.Email,
	}))
}

func (c *UserController) Export() func(string) {
	return export(c, map[string]string{
		"POST:  /modify_pwd": "ModifyPassword",
	})
}
