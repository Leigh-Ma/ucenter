package wb

type QPvpMsg struct {
	Code int32
	Data string /*JSON marshaled data*/
}

type QPvpRound struct{
	Answer string
	IsCorrect bool
	SubjectId string
}