package user

import (
	"github.com/astaxie/beego"
	"ucenter/library/pay"
	"ucenter/models"
	"ucenter/controllers"
	"encoding/xml"
)


type wxResp struct {
	XMLName        xml.Name `xml:"xml"`
	ReturnCode     string   `xml:"return_code"`
	ReturnMsg      string   `xml:"return_msg"`
}


type orderCbController struct {
	beego.Controller
}

func (c *orderCbController) WxCb() {
	code, msg := "FAIL", ""

	r, err := pay.WxParseResult(c.Ctx.Request)
	if err == nil {
		code = "SUCCESS"
		msg = "OK"

		order := models.GetOrder(r.LocalOrderId())
		order.WxNotify(r)
		models.Upsert(order)
	} else {
		msg = err.Error()
	}

	c.Data["xml"] = &wxResp{ReturnCode: code, ReturnMsg: msg}

	c.ServeXML()
}

func (c *orderCbController) AliCb() {
	msg := ""
	c.Ctx.Output.ContentType("xml")

	r, err := pay.AliParseResult(c.Ctx.Request)
	if err == nil {
		msg = "success"

		order := models.GetOrder(r.LocalOrderId())
		order.WxNotify(r)
		models.Upsert(order)
	} else {
		msg = err.Error()
	}

	c.Ctx.Output.Body([]byte(msg))
}

func (c *orderCbController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"GET:  /ali_notify":    "AliCb",
		"GET:  /wx_notify":     "WxCb",
		"POST: /ali_notify":    "AliCb",
		"POST: /wx_notify":     "WxCb",
	})
}
