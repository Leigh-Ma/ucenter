package models

type PvpLog struct {
	TCom
	PlayerId int64
	//OpponentId int64
	PvpId     string
	Level     int
	Round     int
	EscapeAt  int
	Right     int
	EarnCoin  int
	IsPvp     bool
	Questions string
	Brief     string
}

func NewPvpLog(playerId int64) *PvpLog {
	return &PvpLog{
		PlayerId: playerId,
	}
}

func (t *PvpLog) TableName() string {
	return "pvp_logs"
}
