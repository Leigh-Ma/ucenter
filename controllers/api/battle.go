package api

import (
	"ucenter/controllers"
	"ucenter/controllers/proto"
	"ucenter/library/http"
	"ucenter/library/wordsbattle"
	"ucenter/models"
)

type battleController struct {
	controllers.ApiController
}

func (c *battleController) renderJson(resp *http.JResp) {
	c.RenderJson(resp)
}

func (c *battleController) Practice() {
	resp := &http.JResp{}
	ws, err := c.WebSocket()
	if err != nil {
		c.renderJson(resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error()))
		return
	}

	//todo test
	player := &models.Player{Name: "Practice", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pvp := wb.GetAPracticeRoom(1)
	p := wb.NewQPvpPlayer(player, ws)

	err = pvp.Join(p, false)
	if err != nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	c.renderJson(resp.Success())
}

func (c *battleController) VsRobot() {
	resp := &http.JResp{}
	ws, err := c.WebSocket()
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	player := &models.Player{Name: "PVE", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pve := wb.GetAPveRoom(player.PvpLevel())
	p := wb.NewQPvpPlayer(player, ws)

	err = pve.Join(p, true)
	if err != nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	c.renderJson(resp.Success())
}

func (c *battleController) Pvp() {
	resp, f := &http.JResp{}, proto.WB_PvpJoinReq{}
	if !c.CheckInputs(f, resp) {
		return
	}

	ws, err := c.WebSocket()
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	player := &models.Player{Name: "PVP", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pvp := wb.GetAPvpRoom(player.PvpLevel(), f.Mode)
	if pvp == nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	p := wb.NewQPvpPlayer(player, ws)

	err = pvp.Join(p, true)

	c.renderJson(resp.Success())
}

func (c *battleController) Invited() {
	resp, f := &http.JResp{}, &proto.WB_PvpInvitedJoinReq{}
	if !c.CheckInputs(f, resp) {
		return
	}

	ws, err := c.WebSocket()
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	player := &models.Player{Name: "PVE", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pvp := wb.GetShareByGuid(f.Guid)
	if pvp == nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	p := wb.NewQPvpPlayer(player, ws)

	err = pvp.Join(p, true)

	c.renderJson(resp.Success())
}

func (c *battleController) CreateShared() {
	resp, f := &http.JResp{}, &proto.WB_PvpCreateReq{}
	if !c.CheckInputs(f, resp) {
		return
	}

	ws, err := c.WebSocket()
	if err != nil {
		resp.Error(http.ERR_WEB_SOCKET_NEEDED, err.Error())
		c.renderJson(resp)
		return
	}

	player := &models.Player{Name: "PVE", Rank: 1, SubRank: 3, GoldCoin: 20}
	player.Id = 1

	pvp := wb.GetAShareRoom(player.PvpLevel(), f.Mode)

	/*set room config*/
	pvp.C.Difficulty = f.Difficulty
	pvp.C.Mode = f.Mode
	pvp.C.SpawnDuration = f.SpawnDuration
	pvp.C.Subject = f.Subject

	if pvp == nil {
		resp.Error(http.ERR_WB_JOIN_BATTLE_FAILED, err.Error())
		c.renderJson(resp)
		return
	}

	p := wb.NewQPvpPlayer(player, ws)

	err = pvp.Join(p, true)

	c.renderJson(resp.Success())
}

func (c *battleController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"GET:  /practice": "Practice",
		"GET:  /vsrobot":  "VsRobot",
		"GET:  /invited":  "Invited",
		"GET:  /pvp":      "Pvp",
		"GET:  /create":   "CreateShared",
	})
}
