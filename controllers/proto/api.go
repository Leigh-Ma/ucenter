package proto

/* a visitor login with uuid and app key, (ucenter)*/
type FVisitorLogin struct {
	Uuid   string `valid:"Required" json:"uuid"`
	AppKey string `json:"app_key"`
}

/* a player to set it's name after login */
/* @/api/authed/player/set_name */
type FSetPlayerName struct {
	Name string `valid:"" json:"name"`
}

/* get player information, playerid = 0 for self */
/* @/api/authed/player/(me) */
type FGetPlayerInfo struct {
	PlayerId int64 `valid:"" json:"player_id"`
}

/* get a player's wrong words list, playerid = 0 for self */
/* @/api/authed/player/wrong_words */
type FGetPlayerWrongWords struct {
	PlayerId int64 `valid:"" json:"player_id"`
}

/* a player to buy some product */
/* @/api/authed/order/wx_pay */
/* @/api/authed/order/ali_pay */
type FBuyProduct struct {
	ProductId string  `valid:"Email" json:"product_id"`
	Amount    int     `valid:"Min(1)" json:"amount"`
	Price     float32 `valid:"" json:"price"`
}

/* weChat oath login, code is get from weChat*/
type FWxLogin struct {
	Code   string `valid:"Required" json:"code"`
	Uuid   string `valid:"Required" json:"uuid"`
	UserId int64  `json:"user_id"` //optional for new user
}

/* @/api/authed/sign/daily */
type fDailySign struct{}

/* @/api/authed/sign/hour */
type fHourSign struct{}

/* no recording for server, server always save detail log */
/* @/api/authed/battle/create*/
type WB_PvpCreateReq struct {
	Mode          string `json:"mode"`
	Subject       string `json:"subject"`
	Difficulty    string `json:"difficulty"`
	SpawnDuration int64  `json:"spawn_duration"`
}

/* find a pvp room(pvp waiting for another player to start) to join */
/* @/api/authed/battle/pvp */
type WB_PvpJoinReq struct {
	Mode       string `json:"mode"`
	Subject    string `json:"subject"`
	Difficulty string `json:"difficulty"`
}

/* join a pvp room by shared pvp link */
/* @/api/authed/battle/invited */
type WB_PvpInvitedJoinReq struct {
	Guid string `json:"room_id"`
}

/* to start a practice using a fake pvp room */
/* @/api/authed/battle/practice */
type Wb_Practice struct {
	Subject    string `json:"subject"`
	Difficulty string `json:"difficulty"`
}

/* to start a pve using a fake pvp room */
/* @/api/authed/battle/vsrobot */
type Wb_Pve struct {
	Subject    string `json:"subject"`
	Difficulty string `json:"difficulty"`
}
