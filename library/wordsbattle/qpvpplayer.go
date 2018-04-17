package wb

import (
	"encoding/json"
	"errors"
	"time"
)

type qPvpWs interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close()
}

type QPvpPlayer struct {
	PlayerId int64
	Side     int8
	Stakes   int32 // 赌注筹码
	GoldCoin int32
	Stamina  int32

	IsRobot   bool
	Escaped   bool
	EscapedRound int
	LastMsgAt int64

	HintMax  int
	HintUsed int
	SkipMax  int
	SkipUsed int

	Answers   map[int]*qPvpAnswer //roundId
	WS        qPvpWs
	robotEcho chan *QPvpMsg
	pvp       *QPvp
}

func NewQPvpPlayer(coin, stamina int32, playerId int64, ws qPvpWs /*base web socket interfaces*/) *QPvpPlayer {
	return QPvpPlayer{
		PlayerId: playerId,
		GoldCoin: coin,
		Stamina:  stamina,
		Answers:  make(map[int]*qPvpAnswer, 0),
		WS:       ws,
	}
}

func NewQPvpRobot(coin, stamina int32) *QPvpPlayer {
	return QPvpPlayer{
		IsRobot:  true,
		GoldCoin: coin,
		Stamina:  stamina,
		Answers:  make(map[int]*qPvpAnswer, 0),
		WS:       nil,
	}
}

func (t *QPvpPlayer) workAsRealPlayer() error {
	if t.pvp == nil {
		return errors.New("QPvpPlayer pvp message handler not started")
	}

	if t.IsRobot || t.WS == nil {
		return errors.New("QPvpPlayer null websocket or robot can not work as player")
	}

	for {
		_, data, err := t.WS.ReadMessage()
		if err != nil {
			return err
		}

		msg := &QPvpMsg{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			return err
		}

		msg.Side = t.Side

		t.pvp.sendMsg(msg)
	}
}

func (t *QPvpPlayer) workAsRobot() error {
	if t.pvp == nil {
		return errors.New("QPvpPlayer pvp message handler not started")
	}

	if !t.IsRobot && !t.Escaped {
		return errors.New("QPvpPlayer player can not work as robot")
	}

	go func(){
		for {
			msg := <-t.robotEcho
			if msg.Code == pvpNotifyPvpEnd {
				return
			}
			//todo how robot act?
			cpy := &QPvpMsg{Code: msg.Code, TimeStamp: msg.TimeStamp, Side: t.Side, Data: msg.Data}
			t.pvp.sendMsg(cpy)
		}
	}()

	return nil
}

//should not change any information in data msg
func (t *QPvpPlayer) notifyPlayer(msg *QPvpMsg) error {
	if msg == nil || t.IsRobot {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = t.WS.WriteMessage(0, data)
	t.WS.Close()

	return err
}

func (t *QPvpPlayer) notifyPlayerError(err error) error {
	msg := &QPvpMsg{}
	t.prepareMsg(msg, pvpNotifyError, err.Error())

	errText, e := json.Marshal(map[string]string{"error": err.Error()})
	if e == nil {
		msg.Data = errText
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = t.WS.WriteMessage(0, data)
	//t.WS.Close()

	return err
}

//should not change any information in data msg
func (t *QPvpPlayer) notifyRobot(msg *QPvpMsg) error {
	t.robotEcho <- msg
	return nil
}

func (t *QPvpPlayer) prepareRoundAnswer(roundId int) *qPvpAnswer {
	a, ok := t.Answers[roundId]
	if !ok {
		a = &qPvpAnswer{RoundId: roundId, Side: t.Side, AnswerAt: time.Now().Unix()}
		t.Answers[roundId] = a
	}

	return a
}

func (t *QPvpPlayer) prepareMsg(msg *QPvpMsg, code int32, payload interface{}) error {
	data, err := json.Marshal(payload)

	msg.Code = code
	msg.Data = string(data)
	msg.Side = t.Side
	msg.TimeStamp = time.Now().Unix()

	return err
}

func (t *QPvpPlayer) markMsg() {
	t.Escaped = false
	t.LastMsgAt = time.Now().Unix()
}