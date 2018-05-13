package wb

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

func (t *qPvp) _typ() string {
	if !t.IsPractice {
		return "PVP"
	} else {
		return "EXE"
	}
}

func (t *qPvp) logPrefix(msg string) string {
	return fmt.Sprintf("%s[%s] Round[%2d%2d]: %s", t._typ(), t.Guid, t.curRound, t.RoundNum, msg)
}

func (t *qPvp) Info(f string, args ...interface{}) {
	logs.Info(t.logPrefix(fmt.Sprintf(f, args...)))
}

func (t *qPvp) Error(f string, args ...interface{}) {
	logs.Error(t.logPrefix(fmt.Sprintf(f, args...)))
}

func (t *qPvp) Alert(f string, args ...interface{}) {
	logs.Alert(t.logPrefix(fmt.Sprintf(f, args...)))
}

func (t *qPvp) Debug(f string, args ...interface{}) {
	logs.Debug(t.logPrefix(fmt.Sprintf(f, args...)))
}
