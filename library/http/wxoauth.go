package http

import (
	"github.com/astaxie/beego/httplib"
	"ucenter/library/pay"
)

const (
	wxAccessTokenUrl = "https://api.weixin.qq.com/sns/auth"
	//?access_token=ACCESS_TOKEN&openid=OPENID
	wxRefreshTokenUrl = "https://api.weixin.qq.com/sns/oauth2/refresh_token"
	//?appid=APPID&grant_type=refresh_token&refresh_token=REFRESH_TOKEN
	wxGetUserInfoUrl = "https://api.weixin.qq.com/sns/userinfo"
	//?access_token=ACCESS_TOKEN&openid=OPENID

	OAuthChannelWx = "wx"
)

//ACCESS AND REFRESH TOKEN ACK
type wxTokenAck struct {
	Token        string `json:"access_token"`
	ExpireIn     int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`

	//fail
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type wxAuthAck struct {
	Nickname string `json:"nickname"`
	OpenId   string `json:"openid"` //app related id
	Sex      int    `json:"sex"`
	Icon     string `json:"headimgurl"`
	UnionId  string `json:"unionid"` //unique id

	//fail
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WxAuthAck struct {
	Token *wxTokenAck
	User  *wxAuthAck
}

func (t *wxTokenAck) Failed() bool {
	return t == nil || t.ErrCode != 0
}

func (t *wxAuthAck) Failed() bool {
	return t == nil || t.ErrCode != 0
}

func WxOAuthGetAccessToken(code string) (*wxTokenAck, error) {
	req := httplib.Get(wxAccessTokenUrl)
	req.Param("appid", pay.WxCfg.AppId)
	req.Param("secret", pay.WxCfg.AppKey)
	req.Param("code", code)
	req.Param("grant_type", "authorization_code")

	ack := &wxTokenAck{}
	err := req.ToJSON(ack)

	return ack, err
}

func WxOAuthGetUserinfo(accesstoken string, openid string) (*wxAuthAck, error) {
	req := httplib.Get(wxGetUserInfoUrl)
	req.Param("access_token", accesstoken)
	req.Param("openid", openid)

	ack := &wxAuthAck{}
	err := req.ToJSON(ack)

	return ack, err
}

func WxOAuthAuthorize(code string) (*WxAuthAck, error) {
	ack := &WxAuthAck{}

	token, err := WxOAuthGetAccessToken(code)
	if err != nil {
		return ack, err
	}

	ack.Token = token
	if token.Failed() {
		return ack, nil
	}

	user, err := WxOAuthGetUserinfo(token.Token, token.OpenId)
	if err != nil {
		return ack, err
	}

	ack.User = user
	return ack, nil
}

func WxOAuthRefreshToken(refreshToken string) (*wxAuthAck, error) {
	req := httplib.Get(wxRefreshTokenUrl)
	req.Param("appid", pay.WxCfg.AppId)
	req.Param("refresh_token", refreshToken)

	ack := &wxAuthAck{}
	err := req.ToJSON(ack)

	return ack, err
}
