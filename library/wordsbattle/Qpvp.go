package wb

import (
	"time"
	"ucenter/library/types"
)

type QPvp struct {
	Guid           string
	Level          int32
	Subject        string

	CreateAt       int64
	StartThreshold int32 //player num when player started
	RoundNum       int
	IsPvp          bool

	allOffLine bool
	curRound    int
	curQuestion *qPvpQuestion
	questions   map[int]*qPvpQuestion

	err    error
	status int
	cmd    chan *qPvpCmd
	msg    chan *QPvpMsg


	players map[int]*QPvpPlayer
}

func NewQPvp(startThreshold, level int32, round int) *QPvp {
	q := &QPvp{
		Guid:           types.NewGuid().String(),
		Level:          level,
		StartThreshold: startThreshold,
		RoundNum:       round,
		IsPvp:          true,
		CreateAt:       time.Now().Unix(),
		msg:            make(chan *QPvpMsg, 5),
		cmd:            make(chan int),
		players:        make(map[string]*QPvpPlayer, 2),
	}

	q.startCtrl()

	return q
}

func (t *QPvp) Join(p *QPvpPlayer, vsRobot ...bool) error {

	t.sendCmd(&qPvpCmd{Code: pvpCmdJoin, Data: p})

	if len(vsRobot) > 0 && vsRobot[0] {
		t.joinARobot(p.GoldCoin, p.Stamina)
	}

	if p.IsRobot {
		p.workAsRobot()
	} else {
		p.workAsRealPlayer()
	}

	return nil
}

func (t *QPvp) startCtrl() {
	go func() {
		ticker := time.Tick(10 * time.Second)
		for {
			select {
			case msg := <-t.msg:
				t.HandleMsg(msg)
			case cmd := <-t.cmd:
				t.HandleCmd(cmd)
			case <-ticker:
				t.CheckTimeout()
			}

			if t.status == 0 {
				return
			}
		}
	}()
}

func (t *QPvp) sendCmd(cmd *qPvpCmd) {
	t.cmd <- cmd
}

func (t *QPvp) sendMsg(msg *QPvpMsg) {
	t.msg <- msg
}

func (t *QPvp) finished() {
	finishOngoingQPvp(t)
}

func (t *QPvp) lvlDiff(level int32) int32 {
	diff := level - t.Level
	if diff < 0 {
		diff = -diff
	}
	return diff
}

func (t *QPvp) getPlayerBySide(side int) *QPvpPlayer {
	p := t.players[side]
	return p
}

func (t *QPvp) allPlayerBrief() []*qPvpPlayerBrief {
	return nil
}

func (t *QPvp) getRoundAnswers(roundId int) map[int]*qPvpAnswer {
	as := make(map[int]*qPvpAnswer)
	for _, player := range t.players {
		as[player.Side] = player.prepareRoundAnswer(roundId)
	}
	return as
}

func (t *QPvp) isAllOffLine() bool {
	if !t.allOffLine {
		for _, player := range t.players {
			if !player.IsRobot && !player.Escaped {
				return false
			}
		}
	}

	t.allOffLine = true

	return t.allOffLine
}

func (t *QPvp) CheckTimeout() {


	ts := time.Now().Unix()
	t.chkStartTimeOut(ts)
	t.chkRoundTimeOut(ts)

}

func (t *QPvp) chkStartTimeOut(now int64) bool {
	if t.curRound >= 0 {
		return false
	}

	past := now - t.CreateAt
	if  past <= pvpCfgStartTimeOut {
		return false
	}

	var p *QPvpPlayer = nil
	for _, sp := range t.players {
		p = sp
		break
	}

	if p == nil {
		t.errorEnd("No user in this pvp")
		return true
	}

	//todo use robot or end this pvp
	t.joinARobot(p.GoldCoin, p.Stamina)

	return true
}

func (t *QPvp) chkRoundTimeOut(now int64) bool {
	if t.curRound <= 0 {
		return false
	}

	past := now - t.curQuestion.QuestionAt
	if  past <= pvpCfgAnswerTimeout {
		return false
	}

	for _, p := range t.players {
		if p.LastMsgAt - t.curQuestion.QuestionAt <= pvpCfgAnswerTimeout {
			continue
		}
		// todo user action time out
		if !p.IsRobot && !p.Escaped {
			p.workAsRobot()
		}
	}
	return false
}

func (t *QPvp) joinARobot(gold, stamina int32) {
	robot := NewQPvpRobot(gold, stamina)
	t.Join(robot)
}