package api

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"strings"
	. "ucenter/controllers"
	. "ucenter/controllers/form"
	"ucenter/library/http"
	"ucenter/library/tools"
)

var export func(ctrl beego.ControllerInterface, r map[string]string) func(string) = Exportor

var ApiRouter = &RouterGroup{
	Namespace: "/api/",

	Routers: map[string]IExport{
		"login":  &LoginController{},
		"token":  &TokenController{},
		"battle": &BattleController{},
	},
}

type apiController struct {
	beego.Controller
}

func (c *apiController) isJsonReq() bool {
	return (c.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest") ||
		(c.Ctx.Input.Header("Accept") == "application/json")
}

func (c *apiController) parseJsonInput(form interface{}) error {
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

	beego.Info("Prased Form Data:", tools.Stringify(form))

	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		for _, e := range errs {
			return errors.New(e.Error())
		}
	}

	return err
}

func (c *apiController) renderJson(resp *http.JResp) {
	c.Data["json"] = resp
	beego.Info("Server Response: ", resp.ErrorCode, ":", resp.ErrorReason)
	c.ServeJSON()
}

func (c *apiController) checkInputs(f interface{}, resp *http.JResp) bool {
	if err := c.parseJsonInput(f); err != nil {
		resp.Error(http.ERR_PARAMS_ERROR, err.Error())
		c.renderJson(resp)
		return false
	}

	return true
}
