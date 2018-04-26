package models

import (
	"ucenter/library/types"
	"fmt"
	"ucenter/library/pay"
)

const (
	OrderChannelAliPay  = "alipay"
	OrderChannelWxPay   = "wxpay"

	OrderStatusInvalid = "invalid" // third party create but not found in local, record
	OrderStatusCreate = "create" //create, not payed
	OrderStatusConfirm = "confirm" //third party confirm received
	OrderStatusFailed = "failed" //third party confirm received
	OrderStatusCancel = "cancel" //third party confirm received
	OrderStatusDone ="done"//server give player what have been bought
)
type Order struct {
	TCom
	OrderId string
	PlayerId int64
	Channel string
	TransId string  //third party generate
	Product string
	Amount  int
	Price   float32 //total cost
	Status  string
	Sign    string
	Status3 string
}

func NewOrder(playerId int64, amount int, product string, price float32) *Order{
	r := &Order{
		PlayerId:  playerId,
		Product: product,
		Price:   price,
		Amount:  amount,
		OrderId: newOrderId(playerId),
		Status:  OrderStatusCreate,
	}
	return r
}

func (*Order) TableName() string{
	return "orders"
}

func (r *Order) SetChannel(ch string) {
	r.Channel = ch
}

func newOrderId(playerId int64) string{
	//for return check
	return fmt.Sprintf("%s-%s", types.NewGuidString(), playerId)
}

func (r *Order) Brief()string{
	return fmt.Sprintf("%s AMOUNT %d", r.Product, r.Amount)
}

func (r *Order) WxNotify(wxNotify pay.IPayResp) {
	r.thirdPartyNotify(wxNotify, OrderChannelWxPay)
}

func (r *Order) AliNotify(aliNotify pay.IPayResp) {
	r.thirdPartyNotify(aliNotify, OrderChannelAliPay)
}

func (r *Order) thirdPartyNotify(wxResp pay.IPayResp, ch string) {
	if r.IsNew() {
		r.dealNewOrder(wxResp, ch)
		return
	}

	switch r.Status {
	case OrderStatusCreate:
	case OrderStatusConfirm:
		r.dealNormalOrder(wxResp)
	}
}

func (r *Order) dealNewOrder(resp pay.IPayResp, ch string) {
	r.TransId = resp.TransId()
	r.Sign    = resp.Signature()
	r.Status3 = resp.Status()

	r.OrderId = resp.LocalOrderId()
	r.Channel = ch
	r.Price   = resp.Price()

	r.Status  = OrderStatusInvalid
}

func (r *Order) dealNormalOrder(resp pay.IPayResp) {
	r.TransId = resp.TransId()
	r.Sign    = resp.Signature()
	r.Status3 = resp.Status()
	if !resp.Success() {
		r.Status = OrderStatusFailed
		return
	}

	r.Status = OrderStatusConfirm
	player := GetPlayer(r.PlayerId)

	if player.IsNew() && player.Bought( r.Price, r.Product, r.Amount){
		r.Status = OrderStatusDone
	}

	Upsert(player)
}