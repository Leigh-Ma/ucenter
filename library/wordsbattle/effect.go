package wb

func (t *qPvp) takeAnswerEffect(p *qPvpPlayer, a *qPvpAnswer) {
	if t.IsPractice {
		return
	}

	if !a.IsCorrect {
		t._checkAllFailed()
		return
	}

	//My answer is right
	a.Combo = p.Combo
	a.Damage = 0.0
	if t.isNormalMode() {
		t._normalModeEffect(p, a)
	} else {
		t._raceModeEffect(p, a)
	}
}

func (t *qPvp) _normalModeEffect(p *qPvpPlayer, a *qPvpAnswer) {
	//My answer is right
	damage, ratio := float32(1.0), float32(1.0)
	if p.Combo > 2 {
		//I have combo before this time
		ratio = 1.2
	}

	for _, player := range t.players {
		if p == player {
			continue
		}

		if answer, ok := player.Answers[t.curRound]; ok {
			// He has answered before me!
			if answer.IsCorrect {
				//His answers is also right ToT, robot always answer after user
				damage *= 0.8
			} else {
				//His answers is wrong =_=
				damage = 1.3
			}
		} else {
			//I'm the first answered this question
			damage = 1.0
		}

		damage = damage * ratio

		player.HP -= damage
		a.Damage += damage //what's the meaning for multiple player?
	}
}

func (t *qPvp) _raceModeEffect(p *qPvpPlayer, a *qPvpAnswer) {

	//My answer is right
	damage, ratio := float32(1.0), float32(1.0)
	if p.Combo > 2 {
		//I have combo before this time
		ratio = 1.2
	}

	for _, player := range t.players {
		if p == player {
			continue
		}

		if answer, ok := player.Answers[t.curRound]; ok {
			// He has answered before me!
			if answer.IsCorrect {
				//His answers is also right ToT, I can't hurt him
				if !p.IsRobot {
					continue
				}
				//robot always answer after player
			} else {
				//His answers is wrong =_=
				damage = 1.2 //TODO
			}
		} else {
			//I'm the first answered this question
			damage = 1.0
		}

		damage = damage * ratio
		player.HP -= damage

		a.Damage += damage
	}
}

func (t *qPvp) _checkAllFailed() {
	for _, player := range t.players {
		answer, ok := player.Answers[t.curRound]
		if !ok {
			return
		}
		if answer.IsCorrect {
			return
		}
	}

	for _, player := range t.players {
		player.HP -= 0.8
	}
}
