package pay

import (
	"fmt"
	"github.com/imzjy/wxpay"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	WxCfg = &wxpay.WxConfig{
		AppId:         "应用程序Id, 从https://open.weixin.qq.com上可以看得到",
		AppKey:        "API密钥, 在 商户平台->账户设置->API安全 中设置",
		MchId:         "商户号",
		NotifyUrl:     "后台通知地址",
		PlaceOrderUrl: "https://api.mch.weixin.qq.com/pay/unifiedorder",
		QueryOrderUrl: "https://api.mch.weixin.qq.com/pay/orderquery",
		TradeType:     "APP",
	}
	wxPay *wxpay.AppTrans = nil
)

func init() {
	var err error = nil
	wxPay, err = wxpay.NewAppTrans(WxCfg)
	if err != nil {
		panic("wx pay config error")
	}
}

func WxPreOrder(orderId string, price float32, desc, clientIp string) (*wxpay.PaymentRequest, error) {
	prepayId, err := wxPay.Submit(orderId, float64(price), desc, clientIp)
	if err != nil {
		return nil, err
	}

	payRequest := wxPay.NewPaymentRequest(prepayId)
	return &payRequest, nil
}

func WxQuery(transId string) (WxPayResult, error) {
	r, err := wxPay.Query(transId)
	if err != nil {
		return WxPayResult(r), err
	}

	//verity sign of response
	resultInMap := r.ToMap()
	wantSign := wxpay.Sign(resultInMap, wxPay.Config.AppKey)
	gotSign := resultInMap["sign"]
	if wantSign != gotSign {
		err = fmt.Errorf("sign not match, want:%s, got:%s", wantSign, gotSign)
	}

	return WxPayResult(r), err
}

func WxParseResult(req *http.Request) (IPayResp, error) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	r, err := wxpay.ParseQueryOrderResult(body)

	//verity sign of response
	resultInMap := r.ToMap()
	wantSign := wxpay.Sign(resultInMap, wxPay.Config.AppKey)
	gotSign := resultInMap["sign"]
	if wantSign != gotSign {
		err = fmt.Errorf("sign not match, want:%s, got:%s", wantSign, gotSign)
	}

	return (*WxPayResult)(&r), err
}

type WxPayResult wxpay.QueryOrderResult

func (t *WxPayResult) Success() bool {
	return t.ReturnCode == "SUCCESS"
}

func (t *WxPayResult) LocalOrderId() string {
	return t.OrderId
}

func (t *WxPayResult) TransId() string {
	return t.TransactionId
}

func (t *WxPayResult) Price() float32 {
	p, _ := strconv.ParseFloat(t.CashFee, 32)
	return float32(p)
}

func (t *WxPayResult) Signature() string {
	return t.Sign
}

func (t *WxPayResult) Status() string {
	return t.ReturnCode
}
