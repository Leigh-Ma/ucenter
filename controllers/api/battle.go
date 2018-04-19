package api

import(
	"ucenter/library/wordsbattle"
	"ucenter/models"
	"ucenter/library/http"
	"github.com/astaxie/beego/logs"
)

type BattleController struct {
	authorizedController
	wsController
}

func (c *BattleController) VsRobot() {
	resp := &http.JResp{}
	logs.Alert(c.Ctx.Request.Header["Origin"])
	ws, err := c.WebSocket(c.apiController)
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	pvp := wb.GetAPracticeRoom(1)
	p := wb.NewQPvpPlayer(
		&models.Player{Name: "xx", Rank: 1, SubRank: 3, GoldCoin: 20},
		20, 20, ws)

	err = pvp.Join(p, false)
	if err != nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	resp.Success()
	c.renderJson(resp)
}

func (c *BattleController) Export() func(string) {
	return export(c, map[string]string{
		"GET:  /vsrobot": "VsRobot",
	})
}
