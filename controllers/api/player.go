package api

import (
	"strings"
	"ucenter/controllers"
	"ucenter/controllers/proto"
	"ucenter/library/http"
	"ucenter/models"
)

type playerController struct {
	authorizedController
}

func (c *playerController) PlayerInfo() {
	resp, f := &http.JResp{}, &proto.FGetPlayerInfo{}
	if !c.CheckInputs(f, resp) {
		return
	}

	var player *models.Player
	if f.PlayerId == 0 || strings.Contains(c.Ctx.Request.URL.String(), "/me") {
		player = c.currentPlayer()
	} else {
		player = models.GetPlayer(f.PlayerId)
	}

	c.renderJson(resp.Success(&http.D{
		"player": player,
	}))
}

func (c *playerController) SetName() {
	resp, f := &http.JResp{}, &proto.FSetPlayerName{}
	if !c.CheckInputs(f, resp) {
		return
	}

	player := c.currentPlayer()
	player.Name = f.Name //validate
	models.Upsert(player)
	//set name?

	c.renderJson(resp.Success(&http.D{
		"player": player,
	}))
}

func (c *playerController) FailedQuestions() {
	resp, f := &http.JResp{}, &proto.FGetPlayerWrongWords{}
	if !c.CheckInputs(f, resp) {
		return
	}

	var player *models.Player
	if f.PlayerId != 0 {
		player = models.GetPlayer(f.PlayerId)
	} else {
		player = c.currentPlayer()
	}

	words, _, err := models.DBH().MultiQuery(player.QueryCond().And("pass", false),
		&models.AnswerLog{},
		"question_id",
		"last_fail",
		"first_fail",
		"keyword",
	)

	if err != nil {
		c.renderJson(resp.Error(http.ERR_DATA_BASE_ERROR, err.Error()))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"player": player,
		"words":  words,
	}))
}

func (c *playerController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"POST: /set_name":    "SetName",
		"GET:  /wrong_words": "FailedQuestions",
		"GET:  /":            "PlayerInfo",
		"GET:  /me":          "PlayerInfo",
	})
}
