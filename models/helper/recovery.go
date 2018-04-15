package helper

import (
	"time"
)

const (
	recoveryTypeStamina = int8(1)
)

type recoverConfig struct {
	recoverDuration int
	recoverVal      int
	valInit         int
	valMax          int

	typ     int8
	instant bool
}

var recoverCfg = map[int8]*recoverConfig{
	recoveryTypeStamina: {
		recoverDuration: 10,
		recoverVal:      10,
		valInit:         10,
		valMax:          10,
		typ:             recoveryTypeStamina,
		instant:         false,
	},
}

type Recovery struct {
	LastRefreshAt int64
	Typ           int8
	Val           int
}

func NewRecovery(typ int8) *Recovery {
	cfg, ok := recoverCfg[typ]
	if !ok {
		panic("recover useage error")
	}
	return &Recovery{
		LastRefreshAt: time.Now().Unix(),
		Typ:           cfg.typ,
		Val:           cfg.valInit,
	}
}

func (r *Recovery) recover() {
	now := time.Now().Unix()
	cfg := recoverCfg[r.Typ]

	past := int(now - r.LastRefreshAt)
	inc := int(0)

	if cfg.instant {
		increase := float64(cfg.recoverVal*past) / float64(cfg.recoverDuration)
		inc = int(increase)
		if inc > 0 {
			compensate := int(float64(cfg.recoverDuration) * (increase - float64(inc)) / float64(cfg.recoverVal))
			past = past - compensate
		}
	} else if past >= cfg.recoverDuration {
		inc = cfg.recoverVal
		past = cfg.recoverDuration
	}

	if inc > 0 {
		r.Val += inc
		r.LastRefreshAt += int64(past)
		if r.Val > cfg.valMax {

			r.Val = cfg.valMax
		}
	}
}

func (r *Recovery) Value() int {
	r.recover()
	return r.Val
}

func (r *Recovery) Use(use int) bool {
	r.recover()
	if r.Val > use {
		r.Val -= use
		return true
	}
	return false
}
