package proto

type WB_Msg struct {
	Code       int         `json:"code"`
	codeString string      `json:"-"`
	TimeStamp  int64       `json:"time_stamp"`
	PlayerId   int64       `json:"player_id"` /* 0 for server message send to client */
	Data       interface{} `json:"data"`
	payload    []byte      `json:"-"` /* marshaled data for this.Data, cache when broadcast */
}

type Wb_Question struct {
	QuestionId int64  `json:"question_id"`
	Detail     string `json:"detail"`   /* question detail */
	Category   string `json:"category"` /* fill blank(with or without len hint), choose */
}

type Wb_Player struct {
	PlayerId int64  `json:"player_id"`
	Name     string `json:"name"`
	ICon     string `json:"icon"`
	Rank     int    `json:"level"`
	Star     int    `json:"star"`
	IsRobot  bool   `json:"is_robot"`
}

//no recording for server, server always save detail log
type WB_PvpCreateReq struct {
	Mode       string `json:"mode"`
	Subject    string `json:"subject"`
	Difficulty string `json:"difficulty"`
}

//find a pvp room(pvp waiting for another player to start) to join
type WB_PvpJoinReq struct {
	Mode       string `json:"mode"`
	Subject    string `json:"subject"`
	Difficulty string `json:"difficulty"`
}

//find a pvp room(pvp waiting for another player to start) to join
type WB_PvpQuitReq struct {
}

type WB_PvpQuitNotify struct {
	Player *Wb_Player `json:"player"`
}

type WB_PvpAnswerReq struct {
	Round  int    `json:"round"`
	Answer string `json:"round"`
}

type WB_PvpAnswerNotify struct {
	Round   int  `json:"round"` /* in case of timeout answer, do not use question id*/
	IsRight bool `json:"is_right"`
}

type WB_PvpHintReq struct {
	Round int `json:"round"`
}

type WB_PvpHintAck struct {
	Round int    `json:"round"`
	Hint  string `json:"hint"`
}

type WB_PvpStartNotify struct {
	Players  []*Wb_Player `json:"players"`
	Question *Wb_Question `json:"question"`
}

type WB_PvpNextRoundNotify struct {
	Round    int `json:"round"`
	Question *Wb_Question
}

type WB_PvpFinishNotify struct {
}
