package wb

import (
	"errors"
	"time"
	"ucenter/library/types"
)

type qPvp struct {
	Guid           types.IdString
	Level          int
	Subject        string
	CreateAt       int64
	StartThreshold int //player num when player started
	RoundNum       int
	IsPvp          bool
	ticker         int //heart beat times
	allOffLine     bool
	curRound       int
	questions      map[int]*qPvpQuestion
	players        map[int]*qPvpPlayer

	curQuestion *qPvpQuestion `json:"-"`
	err         error         `json:"-"`
	status      int           `json:"-"`
	cmd         chan *qPvpCmd `json:"-"`
	msg         chan *QPvpMsg `json:"-"`
}

func newQPvp(startThreshold, level, round int) *qPvp {
	q := &qPvp{
		Guid:           types.NewIdString(),
		Level:          level,
		StartThreshold: startThreshold,
		RoundNum:       round,
		IsPvp:          true,
		CreateAt:       time.Now().Unix(),
		msg:            make(chan *QPvpMsg, 5),
		cmd:            make(chan *qPvpCmd, 1),
		players:        make(map[int]*qPvpPlayer, 2),
	}

	q.startCtrlRoutine()

	return q
}

func (t *qPvp) Join(p *qPvpPlayer, vsRobot ...bool) error {
	t.sendCmd(&qPvpCmd{Code: pvpCmdJoin, Data: p})
	p.pvp = t

	//if vs a robot or just practice, join a robot as opponent
	if len(vsRobot) > 0 && vsRobot[0] {
		if err := t.joinARobot(p); err != nil {
			t.sendCmd(&qPvpCmd{Code: pvpCmdErrEnd, Data: err})
			return err
		}
	}

	if p.IsRobot {
		//start a robot process routine
		return p.workAsRobot()
	} else {
		//work as real player, just a loop in current routine
		return p.workAsPlayer()
	}

}

func (t *qPvp) CheckTimeout() {
	t.ticker += 1
	t.Debug("Tick[%3d], question: %v", t.ticker, t.curQuestion)
	ts := time.Now().Unix()

	t.chkStartTimeOut(ts)
	t.chkRoundTimeOut(ts)
}

func (t *qPvp) chkStartTimeOut(now int64) bool {
	if t.curRound >= 0 {
		return false
	}

	past := now - t.CreateAt
	if past <= pvpCfgStartTimeOut {
		return false
	}

	var p *qPvpPlayer = nil
	for _, sp := range t.players {
		p = sp
		break
	}

	if p == nil {
		t.errorEnd(errors.New("No user in this pvp"))
		return true
	}

	//todo use robot or end this pvp
	t.joinARobot(p)

	return true
}

func (t *qPvp) chkRoundTimeOut(now int64) bool {
	if t.curRound <= 0 {
		return false
	}

	past := now - t.curQuestion.QuestionAt
	if past <= pvpCfgAnswerTimeout {
		return false
	}

	for _, p := range t.players {
		if now-p.LastMsgAt <= pvpCfgAnswerTimeout {
			continue
		}
		// todo user action time out
		if !p.IsRobot && !p.Escaped {
			t.sendCmd(&qPvpCmd{Code: pvpCmdEscape, Data: p})
		}
	}
	return false
}

func (t *qPvp) joinARobot(mapper *qPvpPlayer) error {
	robot := NewQPvpRobot(mapper.mp, mapper.GoldCoin, mapper.Stamina)
	return t.Join(robot)
}
