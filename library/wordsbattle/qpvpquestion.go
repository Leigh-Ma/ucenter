package wb

import (
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
		p.Combo += 1
	} else {
		p.Combo = 0
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
	if t.IsPractice {
		more := 5
		t._cacheQuestion(t.RoundNum, more)
		t.RoundNum += more
	} else {
		t._cacheQuestion(0, t.RoundNum)
	}
}

func (t *qPvp) getNewQuestion() (q *qPvpQuestion, err error) {
	if t.RoundNum >= t.curRound {
		q = t.questions[t.curRound]
	} else if t.IsPractice {
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
