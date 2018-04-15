package wb

import (
	"ucenter/library/types"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
)

const (
	A = "A"
	B = "B"
)




type QPvp struct {
	Guid string
	Level int32
	RoundNum int32
	AWinRounds   int32
	ctrl chan int

	A *QPvpPlayer
	B *QPvpPlayer
}

func NewQPvp(playerId int64, level int32) *QPvp {
	return &QPvp{
		Guid: types.NewGuid().String(),
		Level: level,
		ctrl: make(chan int),
	}
}



func (t *QPvp) WaitByA(a *QPvpPlayer, vsRobot... bool) error {
	t.A = a

	t.pvp()
	if len(vsRobot) > 0 && vsRobot[0] {
		robot := NewQPvpRobot(a.GoldCoin, a.Stamina)
		robot.SetSide(B, nil)
		t.StartByB(robot)
	}
	return nil
}

func (t *QPvp) StartByB(b *QPvpPlayer) error{

	if t.B != nil {
		err := errors.New("QPvp is restarted when B side is not nil")
		logs.Alert(err.Error())
		if b.PlayerId != t.B.PlayerId {
			return err
		}
	}
	t.B = B

}

func (t *QPvp) Finished() {
	finishOngoingQPvp(t)
}

func (t *QPvp) levelDiff(level int32) int32{
	diff := level - t.Level
	if diff < 0 {
		diff = -diff
	}
	return diff
}

func (t *QPvp) pvp() int32{
	go func() {
		for {
			select {
			case
			}
		}
	}

}



