package wb

import (
	"ucenter/models"
	"encoding/json"
)

func (t *qPvp) doAnswerLog(p *qPvpPlayer, a *qPvpAnswer) {
	if p.mp == nil || p.IsRobot || p.Escaped {
		return
	}

	l := models.GetAnswerLog(p.mp.GetId(), t.curQuestion.QuestionId)

	l.Answer   = a.Answer
	l.PlayerId = p.mp.Id
	l.Hinted   = a.Hinted
	l.PvpId    = t.Guid.ToString()

	if a.IsCorrect {
		l.Right += 1
	} else {
		l.Failed += 1
	}
	
	if l.IsNew() && !a.IsCorrect{
		l.FirstFail = a.AnswerAt
	}

	models.Upsert(l)
}

func (t *qPvp) doPvpLog() {
	qs := []int64{}
	for i, q := range t.questions {
		if i > t.curRound {
			break
		}
		qs = append(qs, q.QuestionId)
	}
	qList := t.asJson(qs)

	for _, p := range t.players {
		if p.IsRobot {
			continue
		}

		l := models.GetPvpLog(p.mp.GetId())
		
		l.PlayerId = p.mp.Id
		l.PvpId =    t.Guid.ToString()
		l.Level =    t.Level
		l.Round =    t.RoundNum
		l.EscapeAt = p.EscapedRound
		l.Right =    p.Right
		l.EarnCoin = int(p.GoldCoin) //TODO
		l.IsPvp =    t.IsPvp && len(t.players) >= 2
		l.Questions = qList
		l.Brief =    t.briefString()

		models.Upsert(l)
	}

}

func (t *qPvp) briefString()string{
	return "brief"
}

func (*qPvp) asJson(data interface{}) string{
	d, _ := json.Marshal(data)
	return string(d)
}