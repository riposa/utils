package mppay

type paymentRequest struct {
	AppId     string
	PartnerId string
	PrepayId  string
	Package   string
	NonceStr  string
	Timestamp string
	Sign      string
}
