package wb

import (
	"encoding/json"
	"errors"
)


func (t *qPvp) handlePlayerCancel(player *qPvpPlayer, msg *QPvpMsg) bool {
	//do not check which player send this cmd
	if !t.IsPractice && t.state >= stateWaiting{
		ack, _ := player.prepareMsg(pvpNotifyCanceled, &qPvpNotifyCancel{Side: player.Side})
		t.broadCastMsg(ack)
		return true
	}

	return false
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
	t.takeAnswerEffect(player, answer)

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
