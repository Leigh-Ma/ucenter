package wb

import (
	"time"
)

type QPvpMsg struct {
	Code      int32
	TimeStamp int
	Side      string /* Server use only, client do not set */
	Data      string /* JSON marshaled data */
}

func (t *QPvp) HandleMsg(msg *QPvpMsg) {
	player := t.getPlayerBySide(msg.Side)
	if player == nil {
		return
	}

	player.markMsg()

	switch msg.Code {
	case pvpMsgAnswerRound:
		t.onMsgAnswerRound(player, msg)
	case pvpMsgAnswerHint:
		t.onMsgAnswerHint(player, msg)
	case pvpMsgAnswerSkip:
		t.onMsgAnswerSkip(player, msg)
	}
}

// robot answer action will be made base on input answer
func (t *QPvp) isAllPlayerAnswered(msg *QPvpMsg) bool {
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

func (t *QPvp) onMsgAnswerRound(player *QPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	answer := t.handlePlayerAnswer(player, msg)
	// NO MORE CHANGE ON msg, please
	if t.IsPvp {
		if t.isAllPlayerAnswered(msg) {
			t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		}
	} else {
		if answer.IsCorrect {
			t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		}
	}
}

func (t *QPvp) onMsgAnswerHint(player *QPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	t.handlePlayerRequestHint(player, msg)
	// NO MORE CHANGE ON msg, please
}

func (t *QPvp) onMsgAnswerSkip(player *QPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	t.handlePlayerSkipRound(player, msg)
	// NO MORE CHANGE ON msg, please
	if t.IsPvp {
		if t.isAllPlayerAnswered(msg) {
			t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		}
	} else {
		t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
	}
}


