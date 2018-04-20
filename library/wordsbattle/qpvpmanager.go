package wb

import (
	"sync"
	"ucenter/library/types"
)

var (
	qPvpWaiting = newQPvpManager()
	qPvpON      = newQPvpManager()
)

func GetAPvpRoom(level int) *qPvp {
	q := qPvpWaiting.matchOneQPvpByLevel(level)
	if q == nil {
		q = newQPvp(2, level, 5)
		qPvpWaiting.addQPvp(q)
	}
	return q
}

func GetAPveRoom(level int) *qPvp {
	q := newQPvp(2, level, 5)
	return q
}

func GetAPracticeRoom(level int) *qPvp {
	q := newQPvp(1, level, 0)
	q.IsPvp = false
	qPvpWaiting.addQPvp(q)
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

func (t *qPvpManager) getQPvp(guid types.IdString) *qPvp {
	t.RLock()
	p := t.PS[guid]
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
