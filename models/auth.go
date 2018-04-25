package models

import (
	"time"
	"ucenter/library/http"
	"ucenter/library/types"
)

const (
	tokenAliveDuration = 24 * 3600
)

type AuthToken struct {
	TCom
	UserId   int64 `orm:"unique"`
	Token    string
	ExpireAt int64
}

func NewAuthToken(userId int64) *AuthToken {
	return &AuthToken{UserId: userId}
}

func (t *AuthToken) TableName() string {
	return "auth_tokens"
}

func (t *AuthToken) VerifyToken(token string) uint {
	var err uint = http.OK

	if token != t.Token {
		err = http.ERR_TOKEN_INVALID
	}

	if t.ExpireAt > time.Now().Unix() {
		err = http.ERR_TOKEN_EXPIRED
	}

	return err
}

func (t *AuthToken) SetNewToken(userId int64, ttl int64) string {

	t.Token = types.RandomString(32)
	t.ExpireAt = time.Now().Unix() + ttl
	t.UserId = userId

	return t.Token
}

func (t *AuthToken) ExpireToken(userId int64) string {
	t.ExpireAt = time.Now().Unix() - 10

	return t.Token
}
