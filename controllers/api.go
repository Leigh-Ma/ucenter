package controllers

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/gorilla/websocket"
	nh "net/http"
	"strings"
	. "ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/library/tools"
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
		(c.Ctx.Input.Header("Accept") == "application/json")
}

func (c *ApiController) parseJsonInput(form interface{}) error {
	var err error = nil

	if strings.Compare(strings.ToUpper(c.Ctx.Request.Method), "GET") == 0 {
		ParseForm(form, c.Input())
	} else {
		err = json.Unmarshal(c.Ctx.Input.RequestBody, form)
	}

	if err != nil {
		beego.Info("JSON request parse err: ", err, c.Ctx.Input.RequestBody, string(c.Ctx.Input.RequestBody))
		return err
	}

	beego.Info(c.Input(), ": Frased Form Data: ", tools.Stringify(form))

	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		for _, e := range errs {
			return errors.New(e.Error())
		}
	}

	return err
}

type IExport interface {
	Export() func(string)
}

type RouterGroup struct {
	Namespace string
	Routers   map[string]IExport
}

func (c *RouterGroup) RegisterRouter() {
	for name, router := range c.Routers {
		router.Export()(c.Namespace + name) //in case of false call to c.Export
	}
}

func Export(ctrl beego.ControllerInterface, r map[string]string) func(string) {
	//"GET: /index" : "Index"

	return func(ctrlNameSpace string) {
		for route, fn := range r {
			ss := strings.SplitN(route, ":", 2)
			match := strings.Trim(ss[1], " ")
			method := strings.ToLower(strings.Trim(ss[0], " "))

			path := ctrlNameSpace + match
			call := method + ":" + fn

			beego.Info(path, " ", call)
			beego.Router(path, ctrl, call)
		}
	}
}
