package wb

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (t *qPvpQuestion) checkAnswer(p *qPvpPlayer, a *qPvpAnswer) {
	a.AnswerAt = time.Now().Unix()
	a.Side = p.Side
	a.RoundId = t.RoundId

	//TODO ANSWER IS RIGHT?
	a.IsCorrect = true
	if a.IsCorrect && !p.Escaped {
		p.Right += 1
	}
}

func (t *qPvp) _cacheQuestion(lastRound, num int) {
	t.Alert("Cacheing questions (%d) + %d", lastRound, num)
	for i := 1; i <= num; i++ {
		t.questions[lastRound+i] = &qPvpQuestion{
			RoundId:    0,
			QuestionId: int64(90000 + i),
			Question:   "question test",
			Hint:       "answer hint",
		}
	}
}

func (t *qPvp) cacheSomeQuestions() {
	if t.questions == nil {
		t.questions = make(map[int]*qPvpQuestion, t.RoundNum)
	}

	//TODO get several questions, numbered by t.RoundNum
	if t.IsPvp {
		t._cacheQuestion(0, t.RoundNum)
	} else {
		more := 5
		t._cacheQuestion(t.RoundNum, more)
		t.RoundNum += more
	}
}

func (t *qPvp) getNewQuestion() (q *qPvpQuestion, err error) {
	if t.RoundNum >= t.curRound {
		q = t.questions[t.curRound]
	} else if !t.IsPvp {
		t.cacheSomeQuestions()
		q = t.questions[t.curRound]
	}

	if q == nil || q.RoundId != 0 {
		err = errors.New(fmt.Sprintf("get question failed, round %d question(%v)", t.curRound, q))
	} else {
		q.RoundId = t.curRound
		q.QuestionAt = time.Now().Unix()
	}

	return
}

func (t *qPvp) handlePlayerAnswer(player *qPvpPlayer, msg *QPvpMsg) *qPvpAnswer {
	answer := player.prepareRoundAnswer(t.curRound)

	err := json.Unmarshal([]byte(msg.Data), answer)
	if err != nil {
		answer.IsCorrect = false
		player.notifyPlayerError(err)
		return answer
	}

	t.curQuestion.checkAnswer(player, answer)
	t.doAnswerLog(player, answer)

	ack, _ := player.prepareMsg(pvpNotifyAnswerCheck, answer)

	player.notifyPlayer(ack)

	return answer
}

func (t *qPvp) handlePlayerRequestHint(player *qPvpPlayer, msg *QPvpMsg) *qPvpHint {
	hint := &qPvpHint{}
	if player.HintUsed >= player.HintMax {
		player.notifyPlayerError(errors.New("Max Hint Time Used"))
		return nil
	}

	err := json.Unmarshal([]byte(msg.Data), hint)
	if err != nil {
		player.notifyPlayerError(err)
		return hint
	}

	if hint.RoundId != t.curRound {
		player.notifyPlayerError(errors.New("Hint round index mismatch with server"))
		return hint
	}

	//set hint
	hint.Hint = t.getHintForPlayer(player)

	ack, _ := player.prepareMsg(pvpNotifyAnswerHint, hint)

	player.notifyPlayer(ack)

	player.HintUsed += 1

	return hint
}

func (t *qPvp) handlePlayerSkipRound(player *qPvpPlayer, msg *QPvpMsg) *qPvpHint {
	skip := &qPvpHint{}

	//if no answer is made, fake one
	player.prepareRoundAnswer(t.curRound)

	if player.SkipUsed >= player.SkipMax {
		player.notifyPlayerError(errors.New("Max Hint Time Used"))
		return skip
	}

	err := json.Unmarshal([]byte(msg.Data), skip)
	if err != nil {
		player.notifyPlayerError(err)
		return skip
	}

	if skip.RoundId != t.curRound {
		player.notifyPlayerError(errors.New("Hint round index mismatch with server"))
		return skip
	}

	skip.Hint = "string" //TODO set it as answer
	if player.IsRobot {
		//TODO robot random set answer
	}

	ack, _ := player.prepareMsg(pvpNotifyAnswerHint, skip)
	//SKIP ALSO GIVE RIGHT ANSWER AS HINT
	player.notifyPlayer(ack)

	player.SkipUsed += 1

	return skip
}
