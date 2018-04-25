package models

func GetUser(userId int64) *User {
	t := NewUser()
	err := t.FindById(t, userId)
	t.isNew = (err == nil)
	return t
}

func GetUserByUuid(uuid string) *User {
	t := NewUser()
	err := t.FindBy("uuid", uuid, t)
	t.isNew = (err == nil)
	return t
}

func GetUserByEmail(email string) *User {
	t := NewUser()
	err := t.FindBy("email", email, t)
	t.isNew = (err == nil)
	return t
}


func GetPlayer(userId int64) *Player{
	t := NewPlayer(userId, "")
	err := t.FindBy("user_id", userId, t)
	t.isNew = (err == nil)
	return t
}

func GetPlayerSign(playerId int64) *PlayerSign {
	t := NewPlayerSign(playerId)
	err := t.FindBy("player_id", playerId, t)
	t.isNew = (err == nil)
	return t
}

func GetAuthToken(userId int64) *AuthToken {
	t := NewAuthToken(userId)
	err := t.FindBy("user_id", userId, t)
	t.isNew = (err == nil)
	return t
}

func GetAnswerLog(playerId, questionId int64) *AnswerLog{
	t := NewAnswerLog(playerId, questionId)
	err := t.NewQuery(t).Filter("player_id", playerId).Filter("question_id", questionId).One(t)
	t.isNew = (err == nil)
	return t
}

func GetPvpLog(playerId int64) *PvpLog {
	t := NewPvpLog(playerId)
	t.isNew = true
	return t
}

func Upsert(obj ITable, cols... string) (int64, error) {
	if obj.IsNew() {
		return obj.Com().Insert(obj)
	} else {
		return obj.Com().Update(obj, cols...)
	}
}

