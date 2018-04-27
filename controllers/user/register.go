package user

import (
	"ucenter/controllers"
	"ucenter/controllers/proto"
	"ucenter/library/http"
	"ucenter/models"
)

type registerController struct {
	controllers.ApiController
}

func (c *registerController) Register() {
	f, resp := &proto.FRegister{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	if f.Password != f.PasswordRe {
		resp.Error(http.ERR_PASSWORD_MISMATCH)
		c.RenderJson(resp)
		return
	}

	user := models.GetUserByEmail(f.Email)
	if !user.IsNew() {
		resp.Error(http.ERR_EMAIL_HAS_BEEN_TAKEN)
		c.RenderJson(resp)
		return
	}

	user.SetPassword(f.Password)
	if _, err := user.Insert(user); err != nil {
		resp.Error(http.ERR_DATA_BASE_ERROR)
		c.RenderJson(resp)
		return
	}

	models.Upsert(user)

	resp.Success(&http.D{"UserId": user.Id, "Email": user.Email})

	c.RenderJson(resp)
}

func (c *registerController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"POST:  /": "Register",
	})
}
