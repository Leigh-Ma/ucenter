package proto

import (
)

type FModifyPassword struct {
	Email       string `valid:"Email"`
	Password    string
	PasswordNew string
}

type FRegister struct {
	Name       string
	Email      string
	Password   string
	PasswordRe string
	Uuid       string
}

type FPhoneRegister struct {
	PhoneID    string `valid:"Phone"`
	VerifyCode string
}

type FTokenVerify struct {
	Token  string
	UserId int64
}

type FTokenLogin struct {
	UserId int64
	Token  string
}

type FPasswordLogin struct {
	Email    string `valid:"Email" json:"email"`
	Password string
}

type FVisitorLogin struct {
	Uuid   string
	AppKey string
}

type FSetPlayerName struct {
	Name string `valid:""`
}

type FBuyProduct struct {
	ProductId string
	Amount    int //>=1
	Price     float32
}
