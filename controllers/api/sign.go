package api

import (
	"ucenter/controllers"
	"ucenter/library/http"
	"ucenter/models"
)

type signController struct {
	authorizedController
}

func (c *signController) Daily() {
	resp := &http.JResp{}

	player := c.currentPlayer()

	sign := models.GetPlayerSign(player.GetId())

	days := sign.DailySign()
	if days == 0 {
		c.renderJson(resp.Error(http.ERR_HAVE_SIGNED_TODAY))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"sign":    sign,
		"rewards": nil,
	}))
}

func (c *signController) Hour() {
	resp := &http.JResp{}
	player := c.currentPlayer()

	sign := models.GetPlayerSign(player.GetId())

	ok := sign.HourSign()
	if !ok {
		c.renderJson(resp.Error(http.ERR_HOUR_SIGN_LATER))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"sign":    sign,
		"rewards": nil,
	}))
}

func (c *signController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"POST:  /daily": "Daily",
		"POST:  /hour":  "Hour",
	})
}
