package wb

import (
	"github.com/astaxie/beego/logs"
)

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
	logs.Info("QPvp begin to handle msg: %s, %+v", msg.codeName(), msg)
	player := t.getPlayerBySide(msg.Side)
	if player == nil {
		logs.Error("QPvp player not found for side", msg.Side)
		return
	}

	switch msg.Code {
	case pvpMsgAnswerRound:
		t.onMsgAnswerRound(player, msg)
	case pvpMsgAnswerHint:
		t.onMsgAnswerHint(player, msg)
	case pvpMsgAnswerSkip:
		t.onMsgAnswerSkip(player, msg)
	}
}

func (t *qPvp) onMsgAnswerRound(player *qPvpPlayer, msg *QPvpMsg) {
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

func (t *qPvp) onMsgAnswerHint(player *qPvpPlayer, msg *QPvpMsg) {
	//wait until other player(s) to finish or broadcast to others ?
	t.handlePlayerRequestHint(player, msg)
	// NO MORE CHANGE ON msg, please
}

func (t *qPvp) onMsgAnswerSkip(player *qPvpPlayer, msg *QPvpMsg) {
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


