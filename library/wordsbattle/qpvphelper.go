package wb

import (
	"time"
	"github.com/astaxie/beego/logs"
)

func (t *qPvp) started() {
	qPvpWaiting.delQPvp(t)
	if t.IsPvp {
		qPvpON.addQPvp(t)
	}

}

func (t *qPvp) finished() {
	qPvpON.delQPvp(t)
	t.status = ctrlStatusFinished
}


func (t *qPvp) moreRound() bool {
	return t.curRound < t.RoundNum
}

func (t *qPvp) sendCmd(cmd *qPvpCmd) {
	if t.status >= ctrlStatusCritical {
		logs.Alert("Cmd send after ctrlStatusCritical: %d, last error %s", cmd.Code, t.err.Error())
		return
	}
	t.cmd <- cmd
}

func (t *qPvp) sendMsg(msg *QPvpMsg) {
	if t.status >= ctrlStatusCritical {
		logs.Alert("Msg send after ctrlStatusCritical: %d, last error %s", msg.Code, t.err.Error())
		return
	}
	t.msg <- msg
}

func (t *qPvp) lvlDiff(level int) int {
	diff := level - t.Level
	if diff < 0 {
		diff = -diff
	}
	return diff
}

// robot answer action will be made base on input answer
func (t *qPvp) isAllPlayerAnswered(msg *QPvpMsg) bool {
	for _, p := range t.players {
		if _, ok := p.Answers[t.curRound]; !ok {
			if p.IsRobot {
				//TODO robot answer question passively according to input answer
				p.notifyRobot(msg)
				// one robot at a time, return
			}
			return false
		}
	}

	//all answered
	if t.curQuestion.AnswerAllAt <= t.curQuestion.QuestionAt {
		t.curQuestion.AnswerAllAt = time.Now().Unix()
	}

	return true
}


func (t *qPvp) getPlayerBySide(side int) *qPvpPlayer {
	p := t.players[side]
	return p
}

func (t *qPvp) allPlayerBrief() (r []*qPvpPlayerBrief){
	for _, p := range t.players {
		r = append(r, p.playerBrief())
	}
	return r
}

func (t *qPvp) getRoundAnswers(roundId int) map[int]*qPvpAnswer {
	as := make(map[int]*qPvpAnswer)
	for _, player := range t.players {
		as[player.Side] = player.prepareRoundAnswer(roundId)
	}
	return as
}

func (t *qPvp) isAllOffLine() bool {
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

func (t *qPvp) canStartPvp() bool {
	return len(t.players) >= int(t.StartThreshold)
}