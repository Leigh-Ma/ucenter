package models

import "ucenter/models/helper"

// Player information for all logic
type Player struct {
	TCom
	UserId int64
	Name string
	PvpLvl int
	GoldCoin int
	helper.Recovery
}

func (t *Player) TableName() string {
	return "players"
}