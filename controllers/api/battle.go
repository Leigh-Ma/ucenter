package api

import (
	"ucenter/library/http"
	"ucenter/library/wordsbattle"
	"ucenter/models"
)

type BattleController struct {
	authorizedController
	wsController
}

func (c *BattleController) Practice() {
	resp := &http.JResp{}
	ws, err := c.WebSocket(c.apiController)
	if err != nil {
		c.renderJson(resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error()))
		return
	}

	player := &models.Player{Name: "Practice", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pvp := wb.GetAPracticeRoom(1)
	p := wb.NewQPvpPlayer(player, 20, 20, ws)

	err = pvp.Join(p, false)
	if err != nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	c.renderJson(resp.Success())
}

func (c *BattleController) VsRobot() {
	resp := &http.JResp{}
	ws, err := c.WebSocket(c.apiController)
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	player := &models.Player{Name: "PVE", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pve := wb.GetAPveRoom(1)
	p := wb.NewQPvpPlayer(player, 20, 20, ws)

	err = pve.Join(p, true)
	if err != nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	c.renderJson(resp.Success())
}

func (c *BattleController) Export() func(string) {
	return export(c, map[string]string{
		"GET:  /practice": "Practice",
		"GET:  /vsrobot":  "VsRobot",
	})
}
