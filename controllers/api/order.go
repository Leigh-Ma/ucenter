package api

import (
	"ucenter/library/http"
	"ucenter/controllers/form"
	"ucenter/models"
	"ucenter/library/pay"
	"ucenter/controllers"
)

type orderController struct {
	authorizedController
}

func (c *orderController) AliPay() {
	resp, f := &http.JResp{}, &form.FBuyProduct{}

	if !c.CheckInputs(f, resp) {
		return
	}

	player := c.currentPlayer()
	order := models.NewOrder(player.GetId(), f.Amount, f.ProductId, float32(f.Amount) * f.Price)
	_, err := order.Insert(order)
	if err != nil {
		c.renderJson(resp.Error(http.ERR_ORDER_PRE_CREATE_ERR, err.Error()))
		return
	}

	url, err := pay.AliPreOrder(order.OrderId, order.Price, order.Brief(), c.currentPlayer().Name)
	if err != nil {
		c.renderJson(resp.Error(http.ERR_ORDER_PRE_CREATE_ERR, err.Error()))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"payurl": url,
	}))
}

func (c *orderController) WxPay() {
	resp, f := &http.JResp{}, &form.FBuyProduct{}

	if !c.CheckInputs(f, resp) {
		return
	}

	player := c.currentPlayer()
	order := models.NewOrder(player.GetId(), f.Amount, f.ProductId,  float32(f.Amount) * f.Price)
	_, err := order.Insert(order)
	if err != nil {
		c.renderJson(resp.Error(http.ERR_ORDER_PRE_CREATE_ERR, err.Error()))
		return
	}

	param, err := pay.WxPreOrder(order.OrderId, order.Price, order.Brief(), c.Ctx.Request.RemoteAddr)
	if err != nil {
		c.renderJson(resp.Error(http.ERR_ORDER_PRE_CREATE_ERR, err.Error()))
		return
	}

	c.renderJson(resp.Success(&http.D{
		"param": param,
	}))
}

func (c *orderController) Export() func(string) {
	return controllers.Export(c, map[string]string{
		"GET: /wx_pay":  "WxPay",
		"GET: /ali_pay": "AliPay",
	})
}