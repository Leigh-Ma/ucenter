package models

import "ucenter/models/helper"

// Player information for all logic
type Player struct {
	TCom
	UserId   int64
	Name     string
	Rank   int
	SubRank int
	Icon string
	GoldCoin int
	helper.Recovery
}

func (t *Player) TableName() string {
	return "players"
}

func (t *Player) DoAnswerLog(questionId string, correct bool) {

}