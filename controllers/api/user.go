package api

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type UserController struct {
	apiController
}

func (c *UserController) ModifyPassword() {
	f, resp := &form.FModifyPassword{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	if f.Password == f.PasswordNew {
		resp.Error(http.ERR_PASSWORD_NOT_CHANGED)
		c.renderJson(resp)
		return
	}

	user, isNew := models.UserM.GetByEmail(f.Email)
	if isNew {
		resp.Error(http.ERR_EMAIL_NOT_REGISTERED)
		c.renderJson(resp)
		return
	}

	user.SetPassword(f.Password)
	if _, err := user.Update(user); err != nil {
		resp.Error(http.ERR_DATA_BASE_ERROR)
		c.renderJson(resp)
		return
	}

	resp.Success(&http.D{"UserId": user.Id, "Email": user.Email})

	c.renderJson(resp)
}
