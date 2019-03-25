package mppay

import "testing"

func TestWx_Refund(t *testing.T) {
	wx := RefundNew(
		"wxcca545d88be2e440",
		"e47ed4ced06742c88de5a9540dc86389",
		"1513715361",
		"https://api.mch.weixin.qq.com/secapi/pay/refund",
	)
	wx.Refund("034eafacf3cb11e888c00242c0a89006",
		100, 100,
		"/Users/hengha/projects/go/src/git.henghajiang.com/backend/dropshipping_mp_order/apiclient_cert.pem",
		"/Users/hengha/projects/go/src/git.henghajiang.com/backend/dropshipping_mp_order/apiclient_key.pem")
}
