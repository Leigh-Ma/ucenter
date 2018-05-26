package pay

import (
	"fmt"
	"github.com/smartwalle/alipay"
	"net/http"
	"strconv"
)

type aliCfg struct {
	AppId        string
	PartnerId    string
	PublicKey    string
	PrivateKey   string
	IsProduction bool
	ReturnUrl    string
	NotifyUrl    string
}

var (
	AliCfg = &aliCfg{}
	//appId, partnerId string, aliPublicKey, privateKey []byte, isProduction bool
	aliPay = alipay.New(AliCfg.AppId,
		AliCfg.PartnerId,
		AliCfg.PublicKey,
		AliCfg.PrivateKey,
		AliCfg.IsProduction)
)

func AliPreOrder(orderId string, price float32, desc, userName string) (string, error) {
	r := alipay.AliPayTradeAppPay{}
	r.NotifyURL = AliCfg.NotifyUrl
	r.ReturnURL = AliCfg.ReturnUrl
	r.Subject = userName + desc
	r.OutTradeNo = orderId
	r.TotalAmount = fmt.Sprintf("%f", price)
	r.ProductCode = desc
	return aliPay.TradeAppPay(r)
}

func AliQuery(orderId string) (*alipay.AliPayTradeQueryResponse, error) {
	r := alipay.AliPayTradeQuery{
		AppAuthToken: string(AliCfg.PublicKey),
		OutTradeNo:   orderId,
	}
	return aliPay.TradeQuery(r)
}

func AliParseResult(req *http.Request) (IPayResp, error) {
	r, err := aliPay.GetTradeNotification(req)

	return (*AliPayResult)(r), err
}

type AliPayResult alipay.TradeNotification

func (t *AliPayResult) Success() bool {
	return t.TradeStatus == "TRADE_SUCCESS" ||
		t.TradeStatus == "TRADE_FINISHED"
}

func (t *AliPayResult) LocalOrderId() string {
	return t.OutTradeNo
}

func (t *AliPayResult) TransId() string {
	return t.TradeNo
}

func (t *AliPayResult) Price() float32 {
	p, _ := strconv.ParseFloat(t.BuyerPayAmount, 32)
	return float32(p)
}

func (t *AliPayResult) Signature() string {
	return t.Sign
}

func (t *AliPayResult) Status() string {
	return t.TradeStatus
}
