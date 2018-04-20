package wb

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"time"
	"ucenter/models"
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

	IsRobot      bool
	Escaped      bool
	EscapedRound int
	LastMsgAt    int64

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

//create robot according to player
func NewQPvpRobot(player *models.Player, coin, stamina int32) *qPvpPlayer {
	return &qPvpPlayer{
		IsRobot:   true,
		GoldCoin:  coin,
		Stamina:   stamina,
		Answers:   make(map[int]*qPvpAnswer, 0),
		WS:        nil,
		mp:        player,
		robotEcho: make(chan *QPvpMsg, 5),
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

	logs.Alert("***********Player[%2d] message reading start*************", t.Side)
	for {
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

	logs.Alert("-------------Player[%2d] message reading end exit [%v]-------------", t.Side, err)

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
		t.WS = nil
	}

	if t.robotEcho == nil {
		t.robotEcho = make(chan *QPvpMsg, 5)
	}

	go func() {
		logs.Alert("***********Robot[%2d] message routine started [%v]*************", t.Side, t.robotEcho == nil)
		for {
			msg := <-t.robotEcho
			logs.Info("Robot[%2d] receive message %s", t.Side, msg.codeName())

			c := msg.Code
			if c == pvpNotifyPvpEnd {
				break
			}

			if c == pvpMsgAnswerRound || c == pvpMsgAnswerHint || c == pvpMsgAnswerSkip {
				msg.Side = t.Side
				t.pvp.sendMsg(msg)
			} else {
				logs.Info("Robot[%2d] ignore echo message %s", t.Side, msg.codeName())
			}
		}
		logs.Alert("-------------Robot[%2d] message routine exit-------------", t.Side)
	}()

	return nil
}

//should not change any information in data msg
func (t *qPvpPlayer) notifyPlayer(msg *QPvpMsg) error {
	if msg == nil || t.IsRobot {
		logs.Error("send player of robot or message nil?")
		return nil
	}

	if msg.Data == "" {
		msg.Data = "{}"
	}

	msg.Cs = msg.codeName()
	data, err := json.Marshal(msg)
	if err != nil {
		logs.Error("send player message[%20s] but marshal error: %s", msg.codeName(), err.Error())
		return err
	}

	err = t.WS.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logs.Error("send player message[%20s] write message error: %s", msg.codeName(), err.Error())
		//todo escape ??
	} else {
		logs.Info("send player[%2d] message: %s", t.Side, msg.codeName())
	}

	return err
}

func (t *qPvpPlayer) notifyPlayerError(err error) error {
	msg, _ := t.prepareMsg(pvpNotifyError, map[string]string{"error": err.Error()})
	logs.Info("notify player[%2d] error: %s, detail: %s", t.Side, msg.codeName(), err.Error())

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
	logs.Info("send robot[%2d] message: %s [%v]", t.Side, msg.codeName(), t.robotEcho == nil)

	t.robotEcho <- msg
	return nil
}

func (t *qPvpPlayer) prepareRoundAnswer(roundId int) *qPvpAnswer {
	a, ok := t.Answers[roundId]
	if !ok {
		a = &qPvpAnswer{RoundId: roundId, Side: t.Side}
		t.Answers[roundId] = a
	}

	return a
}

func (t *qPvpPlayer) prepareMsg(code int32, payload interface{}) (*QPvpMsg, error) {
	msg := &QPvpMsg{Code: code, Side: t.Side, TimeStamp: time.Now().Unix()}

	data, err := json.Marshal(payload)
	if err == nil {
		msg.Data = string(data)
	} else {
		msg.Data = string(err.Error())
	}

	return msg, err
}

func (t *qPvpPlayer) markRecvMsg() {
	t.Escaped = false
	t.LastMsgAt = time.Now().Unix()
}

func (t *qPvpPlayer) playerBrief() *qPvpPlayerBrief {
	if t.IsRobot {
		return &qPvpPlayerBrief{
			Id:      0,
			Side:    t.Side,
			Name:    "ROBOT",
			Rank:    t.mp.Rank,
			SubRank: t.mp.SubRank,
			Icon:    "ICON_ROBOT",
		}
	}

	return &qPvpPlayerBrief{
		Id:      t.mp.Id,
		Side:    t.Side,
		Name:    t.mp.Name,
		Rank:    t.mp.Rank,
		SubRank: t.mp.SubRank,
		Icon:    t.mp.Icon,
	}
}
