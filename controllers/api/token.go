package api

import (
	"ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/models"
)

type TokenController struct {
	apiController
}

// GET /token/verify
func (c *TokenController) Verify() {
	f, resp := &form.FTokenVerify{}, &http.JResp{}

	if !c.checkInputs(f, resp) {
		return
	}

	status := models.TokenM.VerifyToken(f.UserId, f.Token)

	resp.Set(status, &http.D{"Token": f.Token})
	c.renderJson(resp)

	return
}

func (c *TokenController) Export() func(string) {
	return export(c, map[string]string{
		"GET:  /verify":     "Verify",
	})
}