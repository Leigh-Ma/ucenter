package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	oneSubRankStep = 5
)

// Player information for all logic
type Player struct {
	TCom
	UserId   int64 //User
	Name     string
	Rank     int
	SubRank  int
	Icon     string
	GoldCoin int

	Stamina       int
	LastRefreshAt int64

	PvpWin      int64
	PvpLose     int64
	PvpWinGold  int64
	PvpLoseGold int64
	Payed       float32
}

func NewPlayer(accountId int64, name string) *Player {
	p := &Player{
		UserId: accountId,
		Name:   name,
	}

	return p
}

func (t *Player) OnInit() {
	t.GoldCoin = 100
	t.Stamina = 99
	t.Rank = 0
	t.SubRank = 1
}

func (t *Player) TableName() string {
	return "players"
}

func (t *Player) QueryCond() *orm.Condition {
	c := orm.NewCondition()
	return c.And("player_id", t.GetId())
}

func (t *Player) PvpLevel() int {
	return t.Rank*oneSubRankStep + t.SubRank
}

func (t *Player) OnPvpWin() {
	t.SubRank += 1
	if t.SubRank >= oneSubRankStep {
		t.Rank += 1
		t.SubRank = 0
	}
}

func (t *Player) OnPvpLose() {
	if t.SubRank <= 0 {
		if t.Rank <= 0 {
			return
		}
		t.Rank -= 1
	}
	t.SubRank -= 1
}

func (t *Player) StaminaVal() int {
	t.recover()
	return t.Stamina
}

func (t *Player) AddStamina(add int) bool {
	t.recover()
	t.Stamina += add
	return true
}

func (t *Player) UseStamina(use int) bool {
	t.recover()
	if t.Stamina > use {
		t.Stamina -= use
		return true
	}
	return false
}

func (t *Player) recover() {
	now := time.Now().Unix()
	if t.Stamina >= 99 {
		t.LastRefreshAt = now
		return
	}

	past := int(now - t.LastRefreshAt)
	inc := int(0)

	//1 point every second
	increase := float64(1*past) / float64(360)
	inc = int(increase)
	if inc > 0 {
		compensate := int(float64(360) * (increase - float64(inc)) / float64(1))
		past = past - compensate
	}

	if inc > 0 {
		t.Stamina += inc
		if t.Stamina >= 99 {
			t.Stamina = 99
			t.LastRefreshAt = now
		} else {
			t.LastRefreshAt += int64(past)
		}
	}
}
