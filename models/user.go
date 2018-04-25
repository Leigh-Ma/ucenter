package models

import (
	"ucenter/library/types"
)

type User struct {
	TCom
	UserName string `orm:"unique;size(64)"`
	Email    string `orm:"unique;size(64)"`
	Password string
	Salt     string
	Uuid     string
	IsActive bool
	IsForbid bool
	RawPwd   string
}

func NewUser() *User{
	return &User{}
}

func (t *User) TableName() string {
	return "users"
}

func (t *User) md5Pwd(pwd string) string {
	if len(t.Salt) == 0 {
		t.Salt = types.RandomString(32)
	}

	return types.MD5(pwd + t.Salt)
}

func (t *User) SetPassword(pwd string) {
	t.Salt = ""
	t.Password = t.md5Pwd(pwd)
}

func (t *User) VerifyPassword(pwd string) bool {
	return types.MD5(pwd+t.Salt) == t.Password
}
