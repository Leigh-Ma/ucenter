package wb

import (
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	"time"
)

func (t *qPvpQuestion) checkAnswer(answer *qPvpAnswer, isRobot bool) {
	//TODO ANSWER IS RIGHT?
	answer.IsCorrect = true
}

func (t *QPvp) cacheSomeQuestions() {
	if t.IsPvp {
		t.questions = make([]*qPvpQuestion, t.RoundNum+1)
		//TODO get several questions, numbered by t.RoundNum
	} else {
		//cache more questions when single practices
		//TODO get more questions and append to t.question
	}
}

func (t *QPvp) getNewQuestion() (q *qPvpQuestion, err error) {
	if t.RoundNum >= t.curRound {
		q = t.questions[t.curRound]
	} else if !t.IsPvp {
		t.cacheSomeQuestions()
		q = t.questions[t.curRound]
	}

	if q == nil || q.RoundId != 0 {
		err = errors.New(fmt.Sprintf("get question failed, round %s", t.curRound))
	} else {
		q.RoundId = t.curRound
		q.QuestionAt = time.Now().Unix()
	}

	return
}

func (t *QPvp) handlePlayerAnswer(player *QPvpPlayer, msg *QPvpMsg) *qPvpAnswer {
	answer := player.prepareRoundAnswer(t.curRound)

	err := json.Unmarshal(msg.Data, answer)
	if err != nil {
		answer.IsCorrect = false
		player.notifyPlayerError(err)
		return answer
	}

	t.curQuestion.checkAnswer(answer, player.IsRobot)

	err = player.prepareMsg(msg, pvpNotifyAnswerCheck, answer)

	if err != nil {
		player.notifyPlayerError(err)
		return answer
	}

	player.notifyPlayer(msg)

	return answer
}

func (t *QPvp) handlePlayerRequestHint(player *QPvpPlayer, msg *QPvpMsg) *qPvpHint {
	hint := &qPvpHint{}
	if player.HintUsed >= player.HintMax {
		player.notifyPlayerError(errors.New("Max Hint Time Used"))
		return
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

	hint.Hint = "string" //TODO

	err = player.prepareMsg(msg, pvpNotifyAnswerHint, hint)
	if err != nil {
		player.notifyPlayerError(err)
		return hint
	}

	player.notifyPlayer(msg)

	player.HintUsed += 1

	return hint
}

func (t *QPvp) handlePlayerSkipRound(player *QPvpPlayer, msg *QPvpMsg) *qPvpHint {
	skip := &qPvpHint{}

	player.prepareRoundAnswer(t.curRound)

	if player.SkipUsed >= player.SkipMax {
		player.notifyPlayerError(errors.New("Max Hint Time Used"))
		return
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

	err = player.prepareMsg(msg, pvpNotifyAnswerHint, skip)
	if err != nil {
		player.notifyPlayerError(err)
		return skip
	}

	//SKIP ALSO GIVE RIGHT ANSWER AS HINT
	player.notifyPlayer(msg)

	player.SkipUsed += 1

	return skip
}
