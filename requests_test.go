package utils

import (
	"github.com/valyala/fasthttp"
	"testing"
)

type callback struct {
}

func (c *callback) Do(req *fasthttp.Request, resp *fasthttp.Response) interface{} {
	return resp.StatusCode()
}

func TestRequests_Get(t *testing.T) {
	var cb callback

	resp, err := Requests.Get("https://global.henghajiang.com/api/v1/district", map[string]string{"TEXT": "1"}, map[string]string{"Authorization": "Token 1"}, &cb)
	if err != nil {
		reqLogger.Exception(err)
	}
	reqLogger.Info(resp.Status())
	reqLogger.Info(resp.ContentType())
	reqLogger.Info(resp.CallbackOutput)
}

func TestRequests_PostJson(t *testing.T) {
	var cb callback
	var data struct {
		Data struct {
			RegionID          string   `json:"region_id"`
			DimensionSelected []string `json:"dimension_selected"`
		} `json:"data"`
	}
	data.Data.RegionID = "65f94359-dfcd-11e7-8bc5-38c98610113c"
	data.Data.DimensionSelected = []string{}

	resp, err := Requests.PostJson("https://global.henghajiang.com/api/v4/ds/product/78d533c2-94a1-11e8-9e01-0242c0a89003", data, map[string]string{"Authorization": "Token 1"}, &cb)
	if err != nil {
		reqLogger.Exception(err)
	}
	reqLogger.Info(resp.Status())
	reqLogger.Info(resp.ContentType())
	reqLogger.Info(resp.CallbackOutput)
}
