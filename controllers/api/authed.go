package api

import (
	"strconv"
	"ucenter/controllers"
	"ucenter/library/http"
	"ucenter/models"
)

const (
	AuthUserId       = "http_user_id"
	AuthSessionToken = "http_session_token"
)

type authorizedController struct {
	controllers.ApiController
	User      models.User
	AuthToken models.AuthToken
	player    *models.Player
	authed    bool
}

func (c *authorizedController) currentUser() *models.User {
	if c.authed {
		return &c.User
	}
	return nil
}

func (c *authorizedController) currentPlayer() *models.Player {
	if c.player == nil {
		c.player = models.GetPlayerByUserId(c.User.GetId())
		if c.player.IsNew() {
			c.player.OnInit()
			models.Upsert(c.player)
		}
	}

	return c.player
}

func (c *authorizedController) renderJson(resp *http.JResp) {
	c.RenderJson(resp)
}

func (c *authorizedController) xPrepare() {
	//check url prefix?
	c.player = nil
	c.authed = false

	resp := &http.JResp{}
	status := uint(http.OK)

	//do token login
	id, err := strconv.ParseInt(c.Ctx.Input.Header(AuthUserId), 10, 64)
	if err != nil || id <= 0 {
		status = http.ERR_PLEASE_RE_LOGIN
	} else if err = c.User.FindById(id, &c.User); err != nil {
		status = http.ERR_PLEASE_RE_LOGIN
	} else if err = c.AuthToken.FindBy("user_id", c.User.GetId(), &c.AuthToken); err != nil {
		resp.Error(http.ERR_PLEASE_RE_LOGIN)
	} else {
		status = c.AuthToken.VerifyToken(c.Ctx.Input.Header(AuthSessionToken))
	}

	if status != http.OK {
		c.renderJson(resp.Error(status))
		return
	}

	c.authed = true
}
