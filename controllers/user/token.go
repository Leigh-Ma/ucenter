package user

import (
	"ucenter/controllers"
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type tokenController struct {
	controllers.ApiController
}

// GET /token/verify
func (c *tokenController) Verify() {
	f, resp := &form.FTokenVerify{}, &http.JResp{}

	if !c.CheckInputs(f, resp) {
		return
	}

	auth := models.GetAuthToken(f.UserId)

	resp.Status(auth.VerifyToken(f.Token), &http.D{"Token": f.Token})

	c.RenderJson(resp)

	return
}

func (c *tokenController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"GET:  /verify": "Verify",
	})
}
