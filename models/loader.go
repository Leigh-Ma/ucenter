package models

//TODO pre-set value will be erased by find, SO, should reset when find return error

func GetUser(userId int64) *User {
	t := NewUser()
	if userId == 0 {
		t.isNew = true
		return t
	}
	err := t.FindById(userId, t)
	t.Id, t.isNew = userId, (err != nil)
	return t
}

func GetUserByUuid(uuid string) *User {
	t := NewUser()
	err := t.FindBy("uuid", uuid, t)
	t.Uuid, t.isNew = uuid, (err != nil)
	return t
}

func GetUserByEmail(email string) *User {
	t := NewUser()
	err := t.FindBy("email", email, t)
	t.Email, t.isNew = email, (err != nil)
	return t
}

func GetOAuthUserByUserId(userId int64, channel string) *OAuthUser {
	t := NewOAuthUser(userId, channel)
	err := t.NewQuery(t).Filter("user_id", userId).Filter("channel", channel).One(t)
	t.UserId, t.Channel, t.isNew = userId, channel, (err != nil)
	return t
}

func GetOAuthUserByOpenId(openId, channel string) *OAuthUser {
	t := NewOAuthUser(0, channel)
	err := t.NewQuery(t).Filter("open_id", openId).Filter("channel", channel).One(t)
	t.OpenId, t.Channel, t.isNew = openId, channel, (err != nil)
	return t
}

func GetPlayerByUserId(userId int64) *Player {
	t := NewPlayer(userId, "")
	err := t.FindBy("user_id", userId, t)
	t.UserId, t.isNew = userId, (err != nil)
	return t
}

func GetPlayer(playerId int64) *Player {
	t := NewPlayer(0, "")
	err := t.FindById(playerId, t)
	t.Id, t.isNew = playerId, (err != nil)
	return t
}

func GetPlayerSign(playerId int64) *PlayerSign {
	t := NewPlayerSign(playerId)
	err := t.FindBy("player_id", playerId, t)
	t.PlayerId, t.isNew = playerId, (err != nil)
	return t
}

func GetAuthToken(userId int64) *AuthToken {
	t := NewAuthToken(userId)
	err := t.FindBy("user_id", userId, t)
	t.UserId, t.isNew = userId, (err != nil)
	return t
}

func GetItem(playerId int64, itemSn string) *Item {
	t := NewItem(playerId, itemSn)
	err := t.NewQuery(t).Filter("player_id", playerId).Filter("name", itemSn).One(t)
	t.PlayerId, t.Name, t.isNew = playerId, itemSn, (err != nil)
	return t
}

func GetOrder(orderSn string) *Order {
	t := &Order{}
	err := t.FindBy("order_id", orderSn, t)
	t.OrderId, t.isNew = orderSn, (err != nil)
	return t
}

func GetAnswerLog(playerId, questionId int64) *AnswerLog {
	t := NewAnswerLog(playerId, questionId)
	err := t.NewQuery(t).Filter("player_id", playerId).Filter("question_id", questionId).One(t)
	t.PlayerId, t.QuestionId, t.isNew = playerId, questionId, (err != nil)
	return t
}

func GetPvpLog(playerId int64) *PvpLog {
	t := NewPvpLog(playerId)
	t.isNew = true
	return t
}

func Upsert(obj ITable, cols ...string) (i int64, e error) {
	if obj.IDB().IsNew() {
		i, e = obj.IDB().Insert(obj)
		if e == nil {
			obj.IDB().MarkOld()
		}
	} else {
		i, e = obj.IDB().Update(obj, cols...)
	}

	return i, e
}
