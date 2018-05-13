package user

import (
	"errors"
	"ucenter/controllers"
	"ucenter/controllers/proto"
	"ucenter/library/http"
	"ucenter/models"
)

type loginController struct {
	controllers.ApiController
}

// GET /login/visitor
func (c *loginController) VisitorLogin() {
	f, resp := &proto.FVisitorLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	user := models.GetUserByUuid(f.Uuid)
	if user.IsNew() {
		user.Insert(user)
	}

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, 24*3600)

	models.Upsert(auth)

	c.RenderJson(resp.Success(&http.D{
		"token": auth,
	}))
}

// GET /login/token
func (c *loginController) TokenLogin() {
	f, resp := &proto.FTokenLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	player := models.GetPlayerByUserId(f.UserId)
	if player.IsNew() {
		c.RenderJson(resp.Error(http.ERR_USER_ID_INVALID))
		return
	}

	auth := models.GetAuthToken(f.UserId)

	if status := auth.VerifyToken(f.Token); http.OK != status {
		resp.Error(status)
		c.RenderJson(resp.Error(status))
		return
	}

	c.RenderJson(resp.Success(&http.D{
		"token": auth,
	}))
}

// GET /login
func (c *loginController) Login() {
	f, resp := &proto.FPasswordLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	user := models.GetUserByEmail(f.Email)
	if user.IsNew() {
		resp.Error(http.ERR_EMAIL_NOT_REGISTERED)
		c.RenderJson(resp)
		return
	}

	if !user.VerifyPassword(f.Password) {
		resp.Error(http.ERR_PASSWORD_ERROR)
		c.RenderJson(resp)
		return
	}

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, 24*3600)

	models.Upsert(auth)

	c.RenderJson(resp.Success(&http.D{
		"token": auth,
	}))
}

//for old user, must use uuid and old user id as a verify for user
func (c *loginController) WxCodeLogin() {
	f, resp := &proto.FWxLogin{}, &http.JResp{}
	if !c.CheckInputs(f, resp) {
		return
	}

	user := models.GetUserByUuid(f.Uuid)

	if user.IsNew() {
		_, err := user.Insert(user)
		if err != nil {
			resp.Error(http.ERR_DATA_BASE_ERROR, err.Error()) //create new user error
			c.RenderJson(resp)
			return
		}
	} else {
		//for old user, must verify uuid and user id
		if f.UserId != user.GetId() {
			resp.Error(http.ERR_USER_ID_INVALID) //create new user error
			c.RenderJson(resp)
			return
		}
	}

	ack, err := http.WxOAuthAuthorize(f.Code)
	if err == nil {
		if ack.Token.Failed() {
			err = errors.New(ack.Token.ErrMsg)
		} else if ack.User.Failed() {
			err = errors.New(ack.User.ErrMsg)
		}
	}

	if err != nil {
		resp.Error(http.ERR_WX_AUTH_BY_CODE_ERR, err.Error())
		c.RenderJson(resp)
		return
	}

	//just upsert all oauth information, do not care new/old user/oauth-user, whatever
	oa := models.GetOAuthUserByOpenId(ack.User.OpenId, http.OAuthChannelWx)

	oa.Channel = http.OAuthChannelWx
	oa.UserId = user.GetId()

	oa.OpenId = ack.User.OpenId
	oa.IconUrl = ack.User.Icon
	oa.Name = ack.User.Nickname
	oa.Sex = ack.User.Sex

	oa.AccessToken = ack.Token.Token
	oa.RefreshToken = ack.Token.RefreshToken
	oa.Expire = ack.Token.ExpireIn

	models.Upsert(oa)

	auth := models.GetAuthToken(user.GetId())

	auth.SetNewToken(user.Id, oa.Expire)
	auth.SetChannelToken(oa.AccessToken, oa.Channel)

	models.Upsert(auth)

	c.RenderJson(resp.Success(&http.D{
		"token": auth,
	}))
}

func (c *loginController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"POST:  /visitor": "VisitorLogin",
		"POST:  /token":   "TokenLogin",
		"POST:  /wx":      "WxCodeLogin",
		"POST:  /":        "Login",
	})
}
