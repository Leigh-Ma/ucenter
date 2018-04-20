package wb

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	ctrlStatusNormal    = 0
	ctrlStatusTinyError = 1
	ctrlStatusCritical  = 2
	ctrlStatusFinished  = 3
)

type qPvpCmd struct {
	Code int
	Data interface{}
}

func (t *qPvpCmd) codeName() string {
	return codeName[t.Code]
}

func (t *qPvp) startCtrlRoutine() {
	t.status = ctrlStatusNormal
	go func() {
		t.Alert("***********Main routine started*************")
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case msg := <-t.msg:
				t.HandleMsg(msg)
			case cmd := <-t.cmd:
				t.HandleCmd(cmd)
			case <-ticker.C:
				t.CheckTimeout()
			}

			if t.status >= ctrlStatusCritical {
				break
			}
		}
		t.Alert("-------------Main routine exit [%v]-------------", t.err)
	}()
}

func (t *qPvp) HandleCmd(cmd *qPvpCmd) {
	t.Info("Handle cmd: %s", cmd.codeName())
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
	case pvpCmdErrEnd:
		t.status = ctrlStatusCritical
		t.onCmdFinish(cmd.Data)
	}
}

func (t *qPvp) onCmdJoin(p interface{}) {
	t.broadCastMsg(t.genJoinMsg(p))
	if t.canStartPvp() {
		t.sendCmd(&qPvpCmd{Code: pvpCmdStart, Data: nil})
	}
}

func (t *qPvp) onCmdStart(p interface{}) {
	t.broadCastMsg(t.genStartMsg())
	t.started()
}

func (t *qPvp) onCmdEscape(p interface{}) {
	t.broadCastMsg(t.genEscapeMsg(p))

	if t.isAllOffLine() {
		t.errorEnd(errors.New("qPvp all user escaped, no player remain"))
	}
}

func (t *qPvp) onCmdNextRound(p interface{}) {
	if t.moreRound() {
		t.broadCastMsg(t.genNextRoundMsg())
	} else {
		t.sendCmd(&qPvpCmd{Code: pvpCmdFinish})
	}
}

func (t *qPvp) onCmdFinish(p interface{}) {
	t.broadCastMsg(t.genFinishMsg())
	t.finished()
}

func (t *qPvp) genJoinMsg(p interface{}) *QPvpMsg {
	player, ok := p.(*qPvpPlayer)
	if !ok {
		t.errorEnd(errors.New("pvp on cmd join, payload not valid player"))
		return nil
	}

	player.Side = len(t.players) + 1 //start from 1
	t.players[player.Side] = player

	notify := &qPvpNotifyJoin{
		Player: player.playerBrief(),
	}

	data, err := json.Marshal(notify)
	if err != nil {
		t.errorEnd(err)
		return nil
	}

	return &QPvpMsg{
		Code:      pvpNotifyPlayerJoin,
		TimeStamp: time.Now().Unix(),
		Side:      pvpSideServer,
		Data:      string(data),
	}
}

func (t *qPvp) genStartMsg() *QPvpMsg {
	t.curRound = 1
	t.cacheSomeQuestions()

	question, err := t.getNewQuestion()
	if err != nil {
		t.errorEnd(err)
		return nil
	}
	t.curQuestion = question

	notify := &qPvpNotifyStart{
		Question: question,
		Players:  t.allPlayerBrief(),
	}

	data, err := json.Marshal(notify)
	if err != nil {
		t.errorEnd(err)
		return nil
	}

	return &QPvpMsg{
		Code:      pvpNotifyPvpStart,
		TimeStamp: time.Now().Unix(),
		Side:      pvpSideServer,
		Data:      string(data),
	}
}

func (t *qPvp) genNextRoundMsg() *QPvpMsg {
	t.curRound += 1
	question, err := t.getNewQuestion()
	if err != nil {
		t.errorEnd(err)
		return nil
	}

	notify := &qPvpNotifyNextRound{
		Question:     question,
		LastQuestion: t.curQuestion,
		LastAnswers:  t.getRoundAnswers(t.curRound),
	}

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
		Data:      string(data),
	}
}

func (t *qPvp) genEscapeMsg(d interface{}) *QPvpMsg {
	ep, ok := d.(*qPvpPlayer)
	if !ok {
		t.errorEnd(errors.New("pvp on cmd escape, payload not valid player"))
		return nil
	}

	ep.Escaped = true
	ep.EscapedRound = t.curRound

	err := ep.workAsRobot()
	if err != nil {
		t.errorEnd(errors.New("pvp on cmd escape: " + err.Error()))
		return nil
	}

	return &QPvpMsg{Code: pvpNotifyPlayerEscape, TimeStamp: time.Now().Unix(), Side: ep.Side}
}

func (t *qPvp) genFinishMsg() *QPvpMsg {
	return &QPvpMsg{
		Code:      pvpNotifyPvpEnd,
		Side:      pvpSideServer,
		TimeStamp: time.Now().Unix(),
	}
}

func (t *qPvp) broadCastMsg(msg *QPvpMsg, escapeSide ...int) {
	if msg == nil {
		t.Info("Broadcat nil message?")
		return
	}

	t.Info("Broadcat message [%20s]", msg.codeName())

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

func (t *qPvp) errorEnd(err error) {
	t.err = err
	t.Alert("Will ending this pvp, reason: %s", err.Error())
	t.sendCmd(&qPvpCmd{Code: pvpCmdErrEnd, Data: err})
}
