package api

import(
	"ucenter/library/http"
	"ucenter/controllers"
	"ucenter/models"
	"ucenter/controllers/form"
)

type playerController struct {
	authorizedController
}

func (c *playerController) SetName() {
	resp, f := &http.JResp{}, &form.FSetPlayerName{}
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
	resp := http.JResp{}
	player := c.currentPlayer()

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
		"POST:   /set_name": "SetName",
		"GET: /wrong_words": "FailedQuestions",
	})
}
