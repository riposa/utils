package mppay

import (
	"encoding/json"
	//"github.com/satori/go.uuid"
	"strconv"
	//"strings"
	"time"
)

type Wx struct {
	cfg *wxConfig
}

type PrePayResponse struct {
	PrePayID       string `json:"prepay_id"`
	OriginPrePayID string `json:"origin_prepay_id"`
	Sign           string `json:"sign"`
	TimeStamp      string `json:"timeStamp"`
	NonceStr       string `json:"nonceStr"`
	SignType       string `json:"signType"`
}

func New(appId, appKey, mchId, notifyUrl, placeOrderUrl, queryOrderUrl string) *Wx {
	return &Wx{
		cfg: &wxConfig{
			AppId:         appId,
			AppKey:        appKey,
			MchId:         mchId,
			NotifyUrl:     notifyUrl,
			PlaceOrderUrl: placeOrderUrl,
			QueryOrderUrl: queryOrderUrl,
			TradeType:     "JSAPI",
		},
	}
}

func RefundNew(appId, appKey, mchId, refundOrderUrl string) *Wx {
	return &Wx{
		cfg: &wxConfig{
			AppId:          appId,
			AppKey:         appKey,
			MchId:          mchId,
			RefundOrderUrl: refundOrderUrl,
			TradeType:      "JSAPI",
		},
	}
}

func (w *Wx) PrePay(orderID, openId, title string, amount float64, attach ...map[string]interface{}) (*PrePayResponse, error) {

	var attachByte string
	if len(attach) == 0 {
		attachByte = ""
	} else {
		attachByte, err := json.Marshal(attach)
		wxLogger.Infof("attach: %s", string(attachByte))
		if err != nil {
			wxLogger.Exception(err)
			return nil, err
		}
	}

	appTrans, err := newAppTrans(w.cfg)
	if err != nil {
		wxLogger.Exception(err)
		return nil, err
	}
	//uid, _ := uuid.NewV1()
	//orderID := strings.Join(strings.Split(uid.String(), "-"), "")
	wxLogger.Infof("orderID: %s, amount: %f, title: %s, openId: %s, string(attachByte): %s", orderID, amount, title, openId, string(attachByte))
	prepayId, err := appTrans.Submit(orderID, amount, title, openId, string(attachByte))
	if err != nil {
		wxLogger.Exception(err)
		return nil, err
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := NewNonceString()
	sign := Sign(map[string]string{
		"appId":     w.cfg.AppId,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + prepayId,
		"signType":  "MD5",
	}, w.cfg.AppKey)

	return &PrePayResponse{
		PrePayID:       "prepay_id=" + prepayId,
		OriginPrePayID: prepayId,
		Sign:           sign,
		TimeStamp:      timeStamp,
		NonceStr:       nonceStr,
		SignType:       "MD5",
	}, nil
}

func (w *Wx) PrePayWithOrderId(orderId, openId, title string, amount float64, attach ...map[string]interface{}) (*PrePayResponse, error) {

	var attachByte string
	if len(attach) == 0 {
		attachByte = ""
	} else {
		attachByte, err := json.Marshal(attach)
		wxLogger.Infof("attach: %s", string(attachByte))
		if err != nil {
			wxLogger.Exception(err)
			return nil, err
		}
	}

	appTrans, err := newAppTrans(w.cfg)
	if err != nil {
		wxLogger.Exception(err)
		return nil, err
	}
	wxLogger.Infof("orderID: %s, amount: %f, title: %s, openId: %s, string(attachByte): %s", orderId, amount, title, openId, string(attachByte))
	prepayId, err := appTrans.Submit(orderId, amount, title, openId, string(attachByte))
	if err != nil {
		wxLogger.Exception(err)
		return nil, err
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonceStr := NewNonceString()
	sign := Sign(map[string]string{
		"appId":     w.cfg.AppId,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + prepayId,
		"signType":  "MD5",
	}, w.cfg.AppKey)

	return &PrePayResponse{
		PrePayID:       "prepay_id=" + prepayId,
		OriginPrePayID: prepayId,
		Sign:           sign,
		TimeStamp:      timeStamp,
		NonceStr:       nonceStr,
		SignType:       "MD5",
	}, nil
}

func (w *Wx) Refund(orderId string, totalFee int, refundFee int, certFile string, keyFile string) error {
	appTrans := &appTrans{Config: w.cfg}
	wxLogger.Infof("orderID: %s, amount: %d apply refund", orderId, refundFee)
	err := appTrans.Refund(orderId, totalFee, refundFee, certFile, keyFile)
	return err
}
