package api

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type RegisterController struct {
	apiController
}

func (c *RegisterController) Register() {
	f, resp := &form.FRegister{}, &http.JResp{}
	if !c.checkInputs(f, resp) {
		return
	}

	if f.Password != f.PasswordRe {
		resp.Error(http.ERR_PASSWORD_MISMATCH)
		c.renderJson(resp)
		return
	}

	user := models.GetUserByEmail(f.Email)
	if !user.IsNew() {
		resp.Error(http.ERR_EMAIL_HAS_BEEN_TAKEN)
		c.renderJson(resp)
		return
	}

	user.SetPassword(f.Password)
	if _, err := user.Insert(user); err != nil {
		resp.Error(http.ERR_DATA_BASE_ERROR)
		c.renderJson(resp)
		return
	}

	models.Upsert(user)

	resp.Success(&http.D{"UserId": user.Id, "Email": user.Email})

	c.renderJson(resp)
}
