package models

type OAuthUser struct {
	TCom
	UserId  int64
	Channel string

	OpenId       string /*普通用户的标识，对当前开发者帐号唯一*/
	Name         string
	Sex          int
	UnionId      string
	AccessToken  string
	RefreshToken string
	Expire       int64
	IconUrl      string
}

func NewOAuthUser(userId int64, channel string) *OAuthUser {
	return &OAuthUser{
		UserId:  userId,
		Channel: channel,
	}
}

func (*OAuthUser) TableName() string {
	return "oauths"
}
