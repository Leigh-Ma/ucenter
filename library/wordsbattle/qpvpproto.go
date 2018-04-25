package wb

const (
	pvpSideServer = 0

	pvpCmdJoin      = 1
	pvpCmdStart     = 2
	pvpCmdEscape    = 3
	pvpCmdNextRound = 4
	pvpCmdFinish    = 5
	pvpCmdErrEnd    = 6

	pvpNotifyError = 10

	pvpMsgAnswerRound = 101
	pvpMsgAnswerHint  = 102
	pvpMsgAnswerSkip  = 103

	pvpNotifyPvpStart     = 1000
	pvpNotifyRoundCreate  = 1001
	pvpNotifyAnswerCheck  = 1002
	pvpNotifyRoundCheck   = 1003
	pvpNotifyAnswerHint   = 1004
	pvpNotifyPlayerJoin   = 1005
	pvpNotifyPlayerEscape = 1006
	pvpNotifyPvpEnd       = 1007
)

var codeName map[int]string = map[int]string{
	pvpSideServer: "pvpSideServer",

	pvpCmdJoin:      "pvpCmdJoin",
	pvpCmdStart:     "pvpCmdStart",
	pvpCmdEscape:    "pvpCmdEscape",
	pvpCmdNextRound: "pvpCmdNextRound",
	pvpCmdFinish:    "pvpCmdFinish",
	pvpCmdErrEnd:    "pvpCmdErrEnd",

	pvpNotifyError: "pvpNotifyError",

	pvpMsgAnswerRound: "pvpMsgAnswerRound",
	pvpMsgAnswerHint:  "pvpMsgAnswerHint",
	pvpMsgAnswerSkip:  "pvpMsgAnswerSkip",

	pvpNotifyPvpStart:     "pvpNotifyPvpStart",
	pvpNotifyRoundCreate:  "pvpNotifyRoundCreate",
	pvpNotifyAnswerCheck:  "pvpNotifyAnswerCheck",
	pvpNotifyRoundCheck:   "pvpNotifyRoundCheck",
	pvpNotifyAnswerHint:   "pvpNotifyAnswerHint",
	pvpNotifyPlayerJoin:   "pvpNotifyPlayerJoin",
	pvpNotifyPlayerEscape: "pvpNotifyPlayerEscape",
	pvpNotifyPvpEnd:       "pvpNotifyPvpEnd",
}

const (
	pvpCfgStartTimeOut  = 10
	pvpCfgAnswerTimeout = 30
)

//client msg data --> pvpMsgAnswerRound
type qPvpQuestion struct {
	RoundId     int    // server, set after answered
	QuestionId  int64 // server
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
	Hinted    bool  `json:"-"`// server
}

type qPvpHint struct {
	RoundId int    // client
	Hint    string // server
}

type qPvpPlayerBrief struct {
	Id      int64
	Side    int
	Name    string
	Rank    int
	SubRank int
	Icon    string
}

type qPvpNotifyNextRound struct {
	Question     *qPvpQuestion
	LastQuestion *qPvpQuestion
	LastAnswers  map[int]*qPvpAnswer //by side id
}

type qPvpNotifyStart struct {
	Question *qPvpQuestion
	Players  []*qPvpPlayerBrief
}

type qPvpNotifyJoin struct {
	Player *qPvpPlayerBrief
}
