package wb

import (
	"sync"
	"ucenter/controllers/proto"
	"ucenter/library/types"
)

const (
	//TODO
	cfgCoinPvp     = 20
	cfgStaminaPvp  = 5
	cfgRewardRatio = float32(0.8)
)

var (
	_M = struct {
		raceWaiting   *qPvpManager
		raceOn        *qPvpManager
		normalWaiting *qPvpManager
		normalOn      *qPvpManager
	}{
		newQPvpManager(),
		newQPvpManager(),
		newQPvpManager(),
		newQPvpManager(),
	}
)

func _matchOnePvp(level int, mode string) *qPvp {
	if mode == proto.Wb_pvp_mode_normal {
		return _M.normalWaiting.matchOneQPvpByLevel(level)
	} else {
		return _M.raceWaiting.matchOneQPvpByLevel(level)
	}
}

func _waitingPvp(q *qPvp) {
	q.state = stateWaiting
	if q.IsPractice {
		return
	}

	if q.isNormalMode() {
		_M.normalWaiting.addQPvp(q)
	} else {
		_M.raceWaiting.addQPvp(q)
	}
}

func _startPvp(q *qPvp) {
	q.state = stateStarted
	if q.IsPractice {
		return
	}

	if q.isNormalMode() {
		_M.normalWaiting.delQPvp(q)
		_M.normalOn.addQPvp(q)
	} else {
		_M.raceWaiting.delQPvp(q)
		_M.raceOn.addQPvp(q)
	}
}

func _finishPvp(q *qPvp) {
	q.state = stateFinished
	if q.IsPractice {
		return
	}

	if q.isNormalMode() {
		_M.normalOn.delQPvp(q)
	} else {
		_M.raceOn.delQPvp(q)
	}
}

func GetAPvpRoom(level int, mode string) *qPvp {
	q := _matchOnePvp(level, mode)
	if q == nil {
		q = newQPvp(2, level, 5)
		//pvp room config
		q.C.Coin = cfgCoinPvp
		q.C.Stamina = cfgStaminaPvp
		q.C.RewardRatio = cfgRewardRatio
		q.C.Mode = mode
		_waitingPvp(q)
	}
	return q
}

func GetAShareRoom(level int, mode string) *qPvp {
	q := newQPvp(2, level, 5)
	q.C.manualCreate = true
	q.C.Mode = mode
	_waitingPvp(q)
	return q
}

//invited by creator
func GetShareByGuid(guid string) *qPvp {
	q := _M.normalWaiting.getQPvp(guid)
	if q == nil {
		q = _M.raceWaiting.getQPvp(guid)
	}
	return q
}

func GetAPveRoom(level int) *qPvp {
	q := newQPvp(2, level, 5)
	return q
}

func GetAPracticeRoom(level int) *qPvp {
	q := newQPvp(1, level, 0)
	q.IsPractice = true
	return q
}

type qPvpManager struct {
	sync.RWMutex
	PS map[types.IdString]*qPvp
}

func newQPvpManager() *qPvpManager {
	return &qPvpManager{
		PS: make(map[types.IdString]*qPvp, 0),
	}
}

func (t *qPvpManager) getQPvp(guid string) *qPvp {

	t.RLock()
	p := t.PS[types.IdString(guid)]
	t.RUnlock()
	return p
}

func (t *qPvpManager) addQPvp(pvp *qPvp) {
	t.Lock()
	t.PS[pvp.Guid] = pvp
	t.Unlock()
}

func (t *qPvpManager) delQPvp(pvp *qPvp) *qPvp {
	t.Lock()
	delete(t.PS, pvp.Guid)
	t.Unlock()
	return pvp
}

func (t *qPvpManager) matchOneQPvpByLevel(level int) (m *qPvp) {
	diff := 1 << 31

	t.Lock()
	for _, pvp := range t.PS {
		if pvp.IsPractice || pvp.C.manualCreate {
			continue
		}
		if d := pvp.lvlDiff(level); d < diff {
			m = pvp
			diff = d
			if diff == 0 {
				break
			}
		}
	}
	if m != nil {
		delete(t.PS, m.Guid)
	}
	t.Unlock()

	return m
}
