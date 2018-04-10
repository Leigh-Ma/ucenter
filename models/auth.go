package models

import (
	"time"
	"sync"
	"ucenter/library/types"
	"ucenter/library/http"
)

const (
	tokenAliveDuration = 24 * 3600
)

var (
	TokenM = NewAuthTokenManager()
)

type AuthTokenManager struct {
	sync.RWMutex
	UserTokens map[int64]*AuthToken
}

func NewAuthTokenManager() *AuthTokenManager{
	return &AuthTokenManager{
		UserTokens: make(map[int64]*AuthToken, 0),
	}
}

func (t *AuthTokenManager) GetUserToken(userId int64) *AuthToken{
	t.RLock()
	auth, ok := t.UserTokens[userId]
	t.RUnlock()

	if !ok {
		auth = &AuthToken{UserId: userId}
		err := auth.FindBy("UserId", userId, auth)
		if err != nil {
			auth.Insert(auth)
		}
		t.Lock()
		t.UserTokens[userId] = auth
		t.Unlock()
	}

	return auth
}

func (t *AuthTokenManager) AddToken(token *AuthToken) {
	t.Lock()
	t.UserTokens[token.UserId] = token
	t.Unlock()
}

func (t *AuthTokenManager) VerifyToken(userId int64, token string) uint {
	t.Lock()
	auth := t.UserTokens[userId]
	if auth == nil {
		auth = &AuthToken{UserId: userId}
		t.UserTokens[userId] = auth
	}
	t.Unlock()

	return auth.verifyToken(token, userId)
}

type AuthToken struct {
	TCom
	UserId   int64 `orm:"unique"`
	Token    string
	ExpireAt int64
}

func (t *AuthToken) TableName() string {
	return "auth_tokens"
}


func (t *AuthToken) verifyToken(token string, userId ...int64) uint {
	var err uint = http.OK

	if len(userId) > 0 {
		if e := t.FindBy("UserId", userId[0], t); e != nil {
			return http.ERR_USER_ID_INVALID
		}
	}

	if token != t.Token {
		err = http.ERR_TOKEN_INVALID
	}

	if t.ExpireAt > time.Now().Unix() {
		err = http.ERR_TOKEN_EXPIRED
	}

	return err
}

func (t *AuthToken) SetNewToken(userId int64, ttl int64) string {
	err := t.FindBy("UserId", userId, t)
	t.Token = types.RandomString(32)
	t.ExpireAt = time.Now().Unix() + ttl
	t.UserId = userId
	if err == nil {
		t.Update(t)
	} else {
		t.Insert(t)
	}

	return t.Token
}

func (t *AuthToken) expireToken(userId int64) string {
	if err := t.FindBy("UserId", userId, t); err != nil {
		t.ExpireAt = time.Now().Unix() - 10
		t.Update(t, "Token", "ExpireAt")
	}

	return t.Token
}
