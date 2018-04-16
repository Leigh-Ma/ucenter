package wb

import (
	"github.com/astaxie/beego/logs"
	"time"
	"encoding/json"
)

type qPvpCmd struct {
	Code int
	Data interface{}
}

func (t *QPvp) HandleCmd(cmd *qPvpCmd) {
	switch cmd.Code {
	case pvpCmdJoin:
		t.onCmdJoin(cmd.Data)
	case pvpCmdStart:
		t.onCmdStart(cmd.Data)
	case pvpCmdEscape:
		t.onCmdEscape(cmd.Data)
	case pvpCmdNextRound:
		t.onCmdNextRound(cmd.Data)
	case pvpCmdFinish:
		t.onCmdFinish(cmd.Data)
	}
}

func (t *QPvp) onCmdJoin(p interface{}) {
	player, ok := p.(*QPvpPlayer)
	if !ok {
		logs.Alert("pvp on cmd join, payload not valid player")
		return
	}
	player.Side = len(t.players) + 1 //start from 1
	t.players[player.Side] = player
	if player.Side >= t.StartThreshold {
		t.sendCmd(&qPvpCmd{Code: pvpCmdStart, nil})
	}
}


func (t *QPvp) onCmdStart(p interface{}) {
	t.cacheSomeQuestions()
	t.broadCastMsg(t.genStartMsg())
}

func (t *QPvp) onCmdEscape(p interface{}) {
	ep, ok := p.(*QPvpPlayer)
	if !ok {
		logs.Alert("pvp on cmd escape, payload not valid player")
		return
	}

	ep.Escaped = true
	ep.EscapedRound = t.curRound

	msg := &QPvpMsg{Code: pvpNotifyPlayerEscape, TimeStamp: time.Now().Unix(), Side: ep.Side}
	t.broadCastMsg(msg, ep.Side)

	if t.isAllOffLine() {
		t.errorEnd(nil)
	}
}

func (t *QPvp) onCmdNextRound(p interface{}) {
	t.broadCastMsg(t.genNextRoundMsg())
}

func (t *QPvp) onCmdFinish(p interface{}) {
	t.broadCastMsg(t.genFinishMsg())
	t.finished()
}


func (t *QPvp) genStartMsg() *QPvpMsg {
	t.curRound = 1
	question, err := t.getNewQuestion()
	if err != nil {
		t.errorEnd(err)
		return
	}

	notify := &qPvpNotifyStart{
		Question: question,
		Players:  t.allPlayerBrief(),
	}

	data, err := json.Marshal(notify)
	if err != nil {
		t.errorEnd(err)
		return
	}

	return &QPvpMsg{
		Code:      pvpNotifyPvpStart,
		TimeStamp: time.Now().Unix(),
		Side:      pvpSideServer,
		Data:      data,
	}
}

func (t *QPvp) genNextRoundMsg() *QPvpMsg {
	question, err := t.getNewQuestion()
	if err != nil {
		t.errorEnd(err)
		return
	}

	notify := &qPvpNotifyNextRound{
		Question: question,
		LastQuestion:  t.curQuestion,
		LastAnswers: t.getRoundAnswers(t.curRound),
	}

	t.curRound += 1
	t.curQuestion = question

	data, err := json.Marshal(notify)
	if err != nil {
		t.errorEnd(err)
		return nil
	}

	return &QPvpMsg{
		Code:      pvpNotifyRoundCreate,
		TimeStamp: time.Now().Unix(),
		Side:      pvpSideServer,
		Data:      data,
	}
}


func (t *QPvp) genFinishMsg() *QPvpMsg {
	return &QPvpMsg{}
}



func (t *QPvp) broadCastMsg(msg *QPvp, escapeSide ...int) {
	if msg == nil {
		return
	}

	if len(escapeSide) == 0 {
		escapeSide = []int{-1}
	}

	for _, p := range t.players {
		if escapeSide[0] == p.Side {
			continue
		}
		if p.IsRobot || p.Escaped {
			p.notifyRobot(msg)
		} else {
			p.notifyPlayer(msg)
		}
	}
}

func (t *QPvp) errorEnd(err error) {
	t.err = err
	t.sendCmd(&qPvpCmd{Code: pvpCmdFinish, Data: err})
}
