package proto

import ()

/* modify user password */
type FModifyPassword struct {
	Email       string `valid:"Email" json:"email"`
	Password    string `json:"password"`
	PasswordNew string `json:"password_new"`
}

/* do user register */
type FRegister struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	PasswordRe string `json:"password_re"`
	Uuid       string `json:"uuid"`
}

/* register user a phone number */
type FPhoneRegister struct {
	PhoneID    string `valid:"Phone" json:"phone_id"`
	VerifyCode string `json:"code"`
}

/* to verify a user provided token, for inner app service node to verify */
type FTokenVerify struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

/* user login server with a token and user id */
type FTokenLogin struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

/* user login with password and email */
type FPasswordLogin struct {
	Email    string `valid:"Email" json:"email"`
	Password string `json:"password"`
}
