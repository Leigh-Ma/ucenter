package controllers

import (

	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/gorilla/websocket"
	nh "net/http"
	"ucenter/library/http"
	"ucenter/library/tools"
	"encoding/json"
)

var _upg = websocket.Upgrader{}

func init() {
	_upg.CheckOrigin = func(r *nh.Request) bool {
		// allow all connections by default
		return true
	}
}

type ApiController struct {
	beego.Controller
}

func (c *ApiController) WebSocket() (*websocket.Conn, error) {
	return _upg.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
}

func (c *ApiController) CheckInputs(f interface{}, resp *http.JResp) bool {
	if err := c.parseJsonInput(f); err != nil {
		resp.Error(http.ERR_PARAMS_ERROR, err.Error())
		c.RenderJson(resp)
		return false
	}

	return true
}

func (c *ApiController) RenderJson(resp *http.JResp) {
	c.Data["json"] = resp
	beego.Info("Server Response: ", resp.ErrorCode, ":", resp.ErrorReason)
	c.ServeJSON()
}

func (c *ApiController) Prepare() {
	if !c.isJsonReq() {
		resp := &http.JResp{}
		c.RenderJson(resp.Error(http.ERR_SERVE_JSON_ONLY))
		return
	}
}

func (c *ApiController) isJsonReq() bool {
	return (c.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest") ||
		(c.Ctx.Input.Header("Content-Type") == "application/json")
}

func (c *ApiController) parseJsonInput(form interface{}) error {
	var err error = nil
	if c.Ctx.Request.Method == "POST" {
		err = json.Unmarshal(c.Ctx.Input.RequestBody, form)
	}

	if err != nil {
		//err = c.Ctx.Request.ParseForm()
		ParseForm(form, c.Input())
	}

	if err != nil {
		beego.Error("Input ", c.Input(), " Parsed error: ", err.Error())
		return err
	}

	beego.Info("Input ", c.Input(), " Parsed As: ", tools.Stringify(form))

	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		for k, e := range errs {
			return errors.New(k + ": " + e.Error())
		}
	}

	return nil
}
