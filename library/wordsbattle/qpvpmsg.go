package wb

import ()

type QPvpMsg struct {
	Code      int32
	TimeStamp int64
	Side      int    `json:"-"` /* Server use only, client do not set */
	Cs        string /* string for code */
	Data      string /* JSON marshaled data, client */
}

func (t *QPvpMsg) codeName() string {
	return codeName[int(t.Code)]
}

func (t *qPvp) HandleMsg(msg *QPvpMsg) {
	t.Info("Handle msg: %s, %+v", msg.codeName(), msg)
	player := t.getPlayerBySide(msg.Side)
	if player == nil {
		t.Error("Player not found for side %d", msg.Side)
		return
	}

	switch msg.Code {
	case pvpMsgCancel:
		t.onMsgCancel(player, msg)
	case pvpMsgAnswerRound:
		t.onMsgAnswerRound(player, msg)
	case pvpMsgAnswerHint:
		t.onMsgAnswerHint(player, msg)
	case pvpMsgAnswerSkip:
		t.onMsgAnswerSkip(player, msg)
	default:
		t.Error("%s should not be proccessed", msg.codeName())
	}
}

func (t *qPvp) onMsgCancel(player *qPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	canceled := t.handlePlayerCancel(player, msg)
	if canceled {
		t.sendCmd(&qPvpCmd{Code: pvpCmdFinish})
		return
	}
	// NO MORE CHANGE ON msg, please
}

func (t *qPvp) onMsgAnswerRound(player *qPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	answer := t.handlePlayerAnswer(player, msg)
	// NO MORE CHANGE ON msg, please
	if t.IsPractice {
		if answer.IsCorrect {
			t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		}
		return
	}

	if t.isAllPlayerAnswered(msg) {
		t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		return
	}

	//race mode, if no robot, give no opportunity for the others after one right answer
	if !t.isNormalMode() && answer.IsCorrect && !t.HasRobot {
		t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
	}
}

func (t *qPvp) onMsgAnswerHint(player *qPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	t.handlePlayerRequestHint(player, msg)
	// NO MORE CHANGE ON msg, please
}

func (t *qPvp) onMsgAnswerSkip(player *qPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	t.handlePlayerSkipRound(player, msg)
	// NO MORE CHANGE ON msg, please
	if t.IsPractice {
		t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
	} else {
		if t.isAllPlayerAnswered(msg) {
			t.sendCmd(&qPvpCmd{Code: pvpCmdNextRound})
		}
	}
}
