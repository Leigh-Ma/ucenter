package wb

import (
	"sync"
)

var (
	qPvpWaiting = newQPvpManager()
	qPvpON = newQPvpManager()
)

func GetOnGoingQPvpByGuid(guid string) *QPvp{
	return qPvpON.getQPvp(guid)
}

func GetAWaitingQPvp(level int32) *QPvp{
	return qPvpWaiting.matchOneQPvpByLevel(level)
}


func finishOngoingQPvp(pvp *QPvp) *QPvp{
	return qPvpON.delQPvp(pvp)
}

type qPvpManager struct {
	sync.RWMutex
	PS map[string]*QPvp
}

func newQPvpManager() *qPvpManager {
	return &qPvpManager{
		PS: make(map[string]*QPvp, 0),
	}
}

func (t *qPvpManager) getQPvp(guid string) *QPvp {
	t.RLock()
	p := t.PS[guid]
	t.RUnlock()
	return p
}

func (t *qPvpManager) addQPvp(pvp *QPvp)  {
	t.Lock()
	t.PS[pvp.Guid] = pvp
	t.Unlock()
}

func (t *qPvpManager) delQPvp(pvp *QPvp) *QPvp {
	t.Lock()
	delete(t.PS, pvp.Guid)
	t.Unlock()
	return pvp
}

func (t *qPvpManager) matchOneQPvpByLevel(level int32)(m *QPvp) {
	diff := 1<<31

	t.Lock()
	for _, pvp := range t.PS {
		if d := pvp.levelDiff(level); d < diff {
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
