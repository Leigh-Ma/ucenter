package user

import (
	"ucenter/controllers"
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type userController struct {
	controllers.ApiController
}

// POST:  /user/modify_pwd
func (c *userController) ModifyPassword() {
	f, resp := &form.FModifyPassword{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	if f.Password == f.PasswordNew {
		c.RenderJson(resp.Error(http.ERR_PASSWORD_NOT_CHANGED))
		return
	}

	user := models.GetUserByEmail(f.Email)
	if user.IsNew() {
		c.RenderJson(resp.Error(http.ERR_EMAIL_NOT_REGISTERED))
		return
	}

	if ok := user.VerifyPassword(f.Password); !ok {
		c.RenderJson(resp.Error(http.ERR_PASSWORD_ERROR))
		return
	}

	user.SetPassword(f.PasswordNew)
	if _, err := user.Update(user); err != nil {
		c.RenderJson(resp.Error(http.ERR_DATA_BASE_ERROR))
		return
	}

	models.Upsert(user)

	//expire token
	c.RenderJson(resp.Success(&http.D{
		"user_id": user.Id,
		"email": user.Email,
	}))
}

func (c *userController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"POST:  /modify_pwd": "ModifyPassword",
	})
}
