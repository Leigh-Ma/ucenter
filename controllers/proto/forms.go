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

/* a visitor login with uuid and app key, (ucenter)*/
type FVisitorLogin struct {
	Uuid   string `valid:"Required" json:"uuid"`
	AppKey string `json:"app_key"`
}

/* a player to set it's name after login */
type FSetPlayerName struct {
	Name string `valid:"" json:"name"`
}

/* get player information, playerid = 0 for self */
type FGetPlayerInfo struct {
	PlayerId int64 `valid:"" json:"player_id"`
}

/* get a player's wrong words list, playerid = 0 for self */
type FGetPlayerWrongWords struct {
	PlayerId int64 `valid:"" json:"player_id"`
}

/* a player to buy some product */
type FBuyProduct struct {
	ProductId string  `valid:"Email" json:"product_id"`
	Amount    int     `valid:"Min(1)" json:"amount"`
	Price     float32 `valid:"Email" json:"price"`
}

type FWxLogin struct {
	Code   string `valid:"Required" json:"code"`
	Uuid   string `valid:"Required" json:"uuid"`
	UserId int64  `json:"user_id"` //optional for new user
}
