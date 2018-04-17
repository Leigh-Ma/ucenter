package wb

const (
	pvpSideServer = 0

	pvpCmdJoin      = 1
	pvpCmdStart     = 2
	pvpCmdEscape    = 3
	pvpCmdNextRound = 4
	pvpCmdFinish    = 5

	pvpNotifyError = 10

	pvpMsgAnswerRound = 101
	pvpMsgAnswerHint  = 102
	pvpMsgAnswerSkip  = 103

	pvpNotifyPvpStart    = 1000
	pvpNotifyRoundCreate = 1001
	pvpNotifyAnswerCheck = 1002
	pvpNotifyRoundCheck  = 1003
	pvpNotifyAnswerHint  = 1004
	pvpNotifyPlayerEscape = 1005
	pvpNotifyPvpEnd       = 1003
)

const (
	pvpCfgStartTimeOut    = 30
	pvpCfgAnswerTimeout   = 60
)

//client msg data --> pvpMsgAnswerRound
type qPvpQuestion struct {
	RoundId     int    // server
	QuestionId  string // server
	Question    string // server
	QuestionAt  int64  // server
	AnswerAllAt int64  // server
	Hint        string `json:"-"` //server, do not show to client
}

type qPvpAnswer struct {
	RoundId int    // client
	Answer  string // client

	Side      int   // server
	IsCorrect bool  // server
	AnswerAt  int64 // server
}

type qPvpHint struct {
	RoundId int    // client
	Hint    string // server
}

type qPvpPlayerBrief struct {
	Side    int
	Name    string
	Rank    int
	SubRank int
	Icon    string
}

type qPvpNotifyNextRound struct {
	Question     qPvpQuestion
	LastQuestion qPvpQuestion
	LastAnswers  map[int]*qPvpAnswer //by side id
}

type qPvpNotifyStart struct {
	Question qPvpQuestion
	Players  []*qPvpPlayerBrief
}
