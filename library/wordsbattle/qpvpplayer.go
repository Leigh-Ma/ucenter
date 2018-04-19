package wb

import (
	"encoding/json"
	"errors"
	"time"
	"ucenter/models"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

/*basic web socket interfaces*/
type qPvpWs interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type qPvpPlayer struct {
	Side     int
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
	pvp       *qPvp `show:"-"`
	mp        *models.Player
}

func NewQPvpPlayer(player *models.Player, coin, stamina int32, ws qPvpWs) *qPvpPlayer {
	return &qPvpPlayer{
		mp:       player,
		GoldCoin: coin,
		Stamina:  stamina,
		Answers:  make(map[int]*qPvpAnswer, 0),
		WS:       ws,
	}
}

func NewQPvpRobot(coin, stamina int32) *qPvpPlayer {
	return &qPvpPlayer{
		IsRobot:  true,
		GoldCoin: coin,
		Stamina:  stamina,
		Answers:  make(map[int]*qPvpAnswer, 0),
		WS:       nil,
	}
}

func (t *qPvpPlayer) workAsPlayer() (err error) {
	if t.pvp == nil {
		err = errors.New("QPvpPlayer pvp message handler not started")
		return
	}

	if t.IsRobot || t.WS == nil {
		err = errors.New("QPvpPlayer null websocket or robot can not work as player")
		return
	}

	t.markRecvMsg()
	for {
		logs.Debug("Waiting for player[%15s] message.... ", t.mp.Name)
		_, data, e := t.WS.ReadMessage()
		if e != nil {
			err = e
			break
		}

		msg := &QPvpMsg{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			logs.Debug("Unmarshal message error %s err: %s.... ", string(data), err.Error())
			break
		}

		msg.Side = t.Side
		t.markRecvMsg()
		t.pvp.sendMsg(msg)
	}

	logs.Alert("QPvp player message loop exit... %s", err.Error())

	return err
}

func (t *qPvpPlayer) workAsRobot() error {
	if t.pvp == nil {
		return errors.New("QPvpPlayer pvp message handler not started")
	}

	if !t.IsRobot && !t.Escaped {
		return errors.New("QPvpPlayer player can not work as robot")
	}

	if t.WS != nil {
		t.WS.WriteMessage(websocket.CloseMessage, nil)
		t.WS.Close()
		t.WS = nil
	}

	if t.robotEcho == nil {
		t.robotEcho = make(chan *QPvpMsg, 5)
	}

	go func(){
		logs.Info("QPvp player[%2d] working as a robot", t.Side)
		for {
			msg := <-t.robotEcho
			logs.Info("QPvp robot[%2d] receive message %s", t.Side, msg.codeName())
			if msg.Code == pvpNotifyPvpEnd {
				break
			}
			//todo how robot act?
			cpy := &QPvpMsg{Code: msg.Code, TimeStamp: msg.TimeStamp, Side: t.Side, Data: msg.Data}
			t.pvp.sendMsg(cpy)
		}
		logs.Alert("QPvp robot message routine exit...")
	}()

	return nil
}

//should not change any information in data msg
func (t *qPvpPlayer) notifyPlayer(msg *QPvpMsg) error {
	if msg == nil || t.IsRobot {
		logs.Error("notify player of robot or message nil?")
		return nil
	}

	msg.Cs = msg.codeName()
	data, err := json.Marshal(msg)
	if err != nil {
		logs.Error("notify player message[%20s] but marshal error: %s", msg.codeName(), err.Error())
		return err
	}

	logs.Info("notify player[%2d] message: %s", t.Side, msg.codeName())

	err = t.WS.WriteMessage(websocket.TextMessage, data)
	//t.WS.Close()
	if err != nil {
		logs.Error("notify player message[%20s] write message error: %s", msg.codeName(), err.Error())
	}

	return err
}

func (t *qPvpPlayer) notifyPlayerError(err error) error {
	msg := &QPvpMsg{}
	t.prepareMsg(msg, pvpNotifyError, err.Error())

	logs.Info("notify player[%2d] error: %s, detail: %s", t.Side, msg.codeName(), err.Error())

	errText, e := json.Marshal(map[string]string{"error": err.Error()})
	if e == nil {
		msg.Data = string(errText)
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
func (t *qPvpPlayer) notifyRobot(msg *QPvpMsg) error {
	logs.Info("notify robot[%2d] message: %s", t.Side, msg.codeName())
	t.robotEcho <- msg
	return nil
}

func (t *qPvpPlayer) prepareRoundAnswer(roundId int) *qPvpAnswer {
	a, ok := t.Answers[roundId]
	if !ok {
		a = &qPvpAnswer{RoundId: roundId, Side: t.Side, AnswerAt: time.Now().Unix()}
		t.Answers[roundId] = a
	}

	return a
}

func (t *qPvpPlayer) prepareMsg(msg *QPvpMsg, code int32, payload interface{}) error {
	data, err := json.Marshal(payload)

	msg.Code = code
	msg.Data = string(data)
	msg.Side = t.Side
	msg.TimeStamp = time.Now().Unix()

	return err
}

func (t *qPvpPlayer) markRecvMsg() {
	t.Escaped = false
	t.LastMsgAt = time.Now().Unix()
}

func (t *qPvpPlayer) playerBrief() *qPvpPlayerBrief {
	return &qPvpPlayerBrief{
		Side: t.Side,
		Name: t.mp.Name,
		Rank: t.mp.Rank,
		SubRank: t.mp.SubRank,
		Icon: t.mp.Icon,
	}
}