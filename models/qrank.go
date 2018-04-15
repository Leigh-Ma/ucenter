package models

//Qualifying randed game
type QRank struct {
	TCom
	PlayerId int64
	Step int8
	LvlInStep int8
	PvpWin  int64
	PvpLose int64
	PvpWinGold int64
	PvpLoseGold int64
}

func (t *QRank) TableName() string {
	return "q_ranks"
}

func (t *QRank) Win() {

}