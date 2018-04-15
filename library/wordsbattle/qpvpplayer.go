package wb

import (
	"errors"
	ws "github.com/gorilla/websocket"
)

type QPvpPlayer struct {
	PlayerId int64
	Side     string //A,B

	Stakes   int32 // 赌注筹码
	GoldCoin int32
	Stamina  int32

	IsRobot    bool
	Escaped  bool
	EscapedAt int32
	Rounds []*QPvpRound
	WS       ws.Conn
}

func NewQPvpPlayer(playerId int64, coin, stamina int32) *QPvpPlayer{
	return QPvpPlayer{
		GoldCoin: coin,
		Stamina: stamina,
		Rounds: make([]*QPvpRound, 5),
	}
}

func NewQPvpRobot(coin, stamina int32) *QPvpPlayer{
	return QPvpPlayer{
		Side: B, //always B
		GoldCoin: coin,
		Stamina: stamina,
		IsRobot: true,
		Rounds: make([]*QPvpRound, 5),
	}
}

func (t *QPvpPlayer) SetSide(side string/*A/B*/, ws *ws.Conn){
	t.Side = side
	t.WS = ws
}

func (t *QPvpPlayer) doCheck() error{
	if t.Side != A || t.Side != B {
		return errors.New("QPvpPlayer side set error")
	}

	if !t.IsRobot && t.WS == nil {
		return errors.New("QPvpPlayer websocker connection not set")
	}

	return nil
}
