package mppay

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"git.henghajiang.com/backend/golang_utils/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// appTrans is abstract of Transaction handler. With appTrans, we can get prepay id
type appTrans struct {
	Config *wxConfig
}

var (
	wxLogger = log.New()
)

// Initialized the appTrans with specific config
func newAppTrans(cfg *wxConfig) (*appTrans, error) {
	if cfg.AppId == "" ||
		cfg.MchId == "" ||
		cfg.AppKey == "" ||
		cfg.NotifyUrl == "" ||
		cfg.QueryOrderUrl == "" ||
		cfg.PlaceOrderUrl == "" ||
		cfg.TradeType == "" {
		return &appTrans{Config: cfg}, errors.New("config field canot empty string")
	}

	return &appTrans{Config: cfg}, nil
}

// refund request
// totalFee、refundFee 单位分
func (t *appTrans) Refund(orderId string, totalFee int, refundFee int, certFile string, keyFile string) error {
	refundXml := t.signedRefundXmlString(orderId, strconv.Itoa(totalFee), strconv.Itoa(refundFee))
	wxLogger.Infof("request xml: %s", refundXml)
	resp, err := doHttpPostWithSsl(t.Config.RefundOrderUrl, []byte(refundXml), certFile, keyFile)
	if err != nil {
		wxLogger.Exception(err)
		return err
	}
	refundResult, err := ParseRefundOrderResult(resp)
	if err != nil {
		wxLogger.Exception(err)
		return err
	}
	wxLogger.Infof("refund result: %+v", refundResult)
	if refundResult.ReturnCode != "SUCCESS" {
		return fmt.Errorf("return code:%s, return msg:%s", refundResult.ReturnCode, refundResult.ReturnMsg)
	}
	return nil
}

func (t *appTrans) signedRefundXmlString(orderId, totalFee, refundFee string) string {
	param := make(map[string]string)
	param["appid"] = t.Config.AppId
	param["mch_id"] = t.Config.MchId
	param["nonce_str"] = NewNonceString()
	param["out_trade_no"] = orderId
	param["out_refund_no"] = orderId
	param["total_fee"] = totalFee
	param["refund_fee"] = refundFee
	sign := Sign(param, t.Config.AppKey)
	param["sign"] = sign
	return ToXmlString(param)
}

// Submit the order to wx pay and return the prepay id if success,
// Prepay id is used for app to start a payment
// If fail, error is not nil, check error for more information
func (t *appTrans) Submit(orderId string, amount float64, desc string, openId string, attach string) (string, error) {

	odrInXml := t.signedOrderRequestXmlString(orderId, fmt.Sprintf("%.0f", amount), desc, openId, attach)
	wxLogger.Infof("request xml: %s", odrInXml)
	resp, err := doHttpPost(t.Config.PlaceOrderUrl, []byte(odrInXml))
	if err != nil {
		return "", err
	}

	placeOrderResult, err := ParsePlaceOrderResult(resp)
	if err != nil {
		return "", err
	}

	//Verify the sign of response
	resultInMap := placeOrderResult.ToMap()
	wxLogger.Infof("wx resp: %+v", resultInMap)

	if placeOrderResult.ReturnCode != "SUCCESS" {
		return "", fmt.Errorf("return code:%s, return desc:%s", placeOrderResult.ReturnCode, placeOrderResult.ReturnMsg)
	}

	if placeOrderResult.ResultCode != "SUCCESS" {
		return "", fmt.Errorf("resutl code:%s, result desc:%s", placeOrderResult.ErrCode, placeOrderResult.ErrCodeDesc)
	}

	return placeOrderResult.PrepayId, nil
}

func (t *appTrans) newQueryXml(transId string) string {
	param := make(map[string]string)
	param["appid"] = t.Config.AppId
	param["mch_id"] = t.Config.MchId
	param["transaction_id"] = transId
	param["nonce_str"] = NewNonceString()

	sign := Sign(param, t.Config.AppKey)
	param["sign"] = sign

	return ToXmlString(param)
}

// Query the order from weixin pay server by transaction id of weixin pay
func (t *appTrans) Query(transId string) (QueryOrderResult, error) {
	queryOrderResult := QueryOrderResult{}

	queryXml := t.newQueryXml(transId)
	// fmt.Println(queryXml)
	resp, err := doHttpPost(t.Config.QueryOrderUrl, []byte(queryXml))
	if err != nil {
		return queryOrderResult, nil
	}

	queryOrderResult, err = ParseQueryOrderResult(resp)
	if err != nil {
		return queryOrderResult, err
	}

	//verity sign of response
	resultInMap := queryOrderResult.ToMap()
	wantSign := Sign(resultInMap, t.Config.AppKey)
	gotSign := resultInMap["sign"]
	if wantSign != gotSign {
		return queryOrderResult, fmt.Errorf("sign not match, want:%s, got:%s", wantSign, gotSign)
	}

	return queryOrderResult, nil
}

// NewPaymentRequest build the payment request structure for app to start a payment.
// Return stuct of paymentRequest, please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=9_12&index=2
func (t *appTrans) NewPaymentRequest(prepayId string) paymentRequest {
	noncestr := NewNonceString()
	timestamp := NewTimestampString()

	param := make(map[string]string)
	param["appid"] = t.Config.AppId
	param["partnerid"] = t.Config.MchId
	param["prepayid"] = prepayId
	param["package"] = "Sign=WXPay"
	param["noncestr"] = noncestr
	param["timestamp"] = timestamp

	sign := Sign(param, t.Config.AppKey)

	payRequest := paymentRequest{
		AppId:     t.Config.AppId,
		PartnerId: t.Config.MchId,
		PrepayId:  prepayId,
		Package:   "Sign=WXPay",
		NonceStr:  noncestr,
		Timestamp: timestamp,
		Sign:      sign,
	}

	return payRequest
}

func (t *appTrans) newOrderRequest(orderId, amount, desc, openId, attach string) map[string]string {
	param := make(map[string]string)
	param["appid"] = t.Config.AppId
	param["attach"] = attach
	param["body"] = desc
	param["mch_id"] = t.Config.MchId
	param["nonce_str"] = NewNonceString()
	param["notify_url"] = t.Config.NotifyUrl
	param["out_trade_no"] = orderId
	param["openid"] = openId
	param["total_fee"] = amount
	param["trade_type"] = t.Config.TradeType
	param["time_start"] = time.Now().Format("20060102150405")
	param["time_expire"] = time.Now().Add(time.Minute * 45).Format("20060102150405")

	return param
}

func (t *appTrans) signedOrderRequestXmlString(orderId, amount, desc, openId, attach string) string {
	order := t.newOrderRequest(orderId, amount, desc, openId, attach)
	sign := Sign(order, t.Config.AppKey)
	// fmt.Println(sign)

	order["sign"] = sign

	return ToXmlString(order)
}

// doRequest post the order in xml format with a sign
func doHttpPost(targetUrl string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}

func doHttpPostWithSsl(targetUrl string, body []byte, certFile string, keyFile string) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return []byte(""), err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil

}
