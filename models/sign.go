package models

import "time"

type PlayerSign struct {
	TCom         `json:"-"`
	PlayerId     int64
	LastSignAt   int64
	SignDays     int
	HourRewardAt int64
}

func NewPlayerSign(playerId int64) *PlayerSign {
	return &PlayerSign{
		PlayerId: playerId,
	}
}

func (t *PlayerSign) TableName() string {
	return "dailys"
}

func (t *PlayerSign) DailySign() int {
	now := time.Now()
	n := now.YearDay()
	o := time.Unix(t.LastSignAt, 0).YearDay()
	if n == o {
		return 0
	} else if n-o == 1 {
		t.LastSignAt = now.Unix()
		t.SignDays += 1
	} else if n-o > 1 {
		t.LastSignAt = now.Unix()
		t.SignDays = 1
	}

	return t.SignDays
}

func (t *PlayerSign) HourSign() bool {
	now := time.Now().Unix()
	if now-t.HourRewardAt >= 3600 {
		t.HourRewardAt = now
		return true
	}
	return false
}
