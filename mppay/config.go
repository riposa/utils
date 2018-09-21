package mppay

type wxConfig struct {
	AppId          string
	AppKey         string
	MchId          string
	NotifyUrl      string
	PlaceOrderUrl  string
	QueryOrderUrl  string
	RefundOrderUrl string
	TradeType      string
}
