package wb

import ()

type pvpAnswer struct {
	RoundIdx   int
	AnswerTime int
	UseHint    bool
}

type pvpQuestion struct {
	QuestionId int64
	Answer     string
	Detail     string
	RoundIdx   int
	Hint       int
}

type pvpPlayer struct {
	Answers []*pvpAnswer
}

type Pvp struct {
	A *pvpPlayer
	B *pvpPlayer

	Questions []*pvpQuestion

	CurRound  int
	MaxRound  int
	RoundStep int
	Status    int
	Level     int
}
