package wb

import (
	"github.com/gorilla/websocket"
	"time"
)

func (t *qPvp) started() {
	qPvpWaiting.delQPvp(t)
	if t.IsPvp {
		qPvpON.addQPvp(t)
	}

}

func (t *qPvp) finished() {
	time.Sleep(2 * time.Second)
	qPvpON.delQPvp(t)
	t.status = ctrlStatusFinished
	for _, player := range t.players {
		if player.WS != nil {
			player.WS.WriteMessage(websocket.CloseMessage, nil)
			player.WS = nil
		}
	}
	//TODO REWARDS?
	t.doPvpLog()
	t.Alert("@@@@@@@@@Over@@@@@@@@@@")
}

func (t *qPvp) moreRound() bool {
	return t.curRound < t.RoundNum
}

func (t *qPvp) sendCmd(cmd *qPvpCmd) {
	if t.status >= ctrlStatusCritical {
		t.Alert("Cmd send after ctrlStatusCritical: %s, last error %s", cmd.codeName(), t.err.Error())
		return
	}
	t.cmd <- cmd
}

func (t *qPvp) sendMsg(msg *QPvpMsg) {
	if t.status >= ctrlStatusCritical {
		t.Alert("Msg send after ctrlStatusCritical: %s, last error %s", msg.codeName(), t.err.Error())
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

func (t *qPvp) allPlayerBrief() (r []*qPvpPlayerBrief) {
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

func (t *qPvp) getHintForPlayer(player *qPvpPlayer) string {
	a := player.prepareRoundAnswer(t.curRound)
	a.Hinted = true
	return t.curQuestion.Hint
}
