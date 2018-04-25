package api

import (
	"ucenter/library/http"
	"ucenter/models"
)

type SignController struct {
	authorizedController
}

func (c *SignController) Daily() {
	resp := &http.JResp{}
	player := c.currentPlayer()
	if player == nil {
		c.renderJson(resp.Error(http.ERR_USER_ID_INVALID))
		return
	}

	sign := models.GetPlayerSign(player.GetId())

	days := sign.DailySign()
	if days == 0 {
		c.renderJson(resp.Error(http.ERR_HAVE_SIGNED_TODAY))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"sign": sign,
		"rewards": nil,
		}))
}

func (c *SignController) Hour() {
	resp := &http.JResp{}
	player := c.currentPlayer()
	if player == nil {
		c.renderJson(resp.Error(http.ERR_USER_ID_INVALID))
		return
	}

	sign := models.GetPlayerSign(player.GetId())

	ok := sign.HourSign()
	if !ok {
		c.renderJson(resp.Error(http.ERR_HOUR_SIGN_LATER))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"sign": sign,
		"rewards": nil,
	}))
}

func (c *SignController) Export() func(string) {
	return export(c, map[string]string{
		"POST:  /daily":    "Daily",
		"POST:  /hour":     "Hour",
	})
}
