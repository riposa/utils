package utils

import (
	"encoding/json"
	"github.com/riposa/utils/log"
	"github.com/valyala/fasthttp"
	url2 "net/url"
	"time"
)

type requests struct {
	// nothing
}

type HTTPCallback interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) interface{}
}

type HTTPResponse struct {
	status      int
	contentType []byte
	body        []byte

	CallbackOutput []interface{}
}

func (h *HTTPResponse) Status() int {
	return h.status
}

func (h *HTTPResponse) ContentType() string {
	return string(h.contentType)
}

func (h *HTTPResponse) Body() []byte {
	return h.body
}

func (h *HTTPResponse) Json(v interface{}) error {
	return json.Unmarshal(h.body, v)
}

func (h *HTTPResponse) Text() string {
	return string(h.body)
}

var (
	Requests  requests
	reqLogger = log.New()
)

func (r *requests) Get(url string, param map[string]string, headers map[string]string, callback ...HTTPCallback) (*HTTPResponse, error) {
	var result HTTPResponse
	var t1, t2, t3, t4 int64
	var cbCount int

	t1 = time.Now().UnixNano()
	u, err := url2.Parse(url)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	reqUri := fasthttp.AcquireURI()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseURI(reqUri)

	reqUri.SetHost(u.Host)
	reqUri.SetPath(u.Path)
	reqUri.SetScheme(u.Scheme)

	queryValues := make(url2.Values)
	for k, v := range param {
		queryValues[k] = []string{v}
	}
	queryString := queryValues.Encode()
	reqUri.SetQueryString(queryString)
	req.SetRequestURIBytes(reqUri.FullURI())
	req.Header.SetMethod("GET")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	t2 = time.Now().UnixNano()
	err = fasthttp.Do(req, resp)
	if err != nil {
		reqLogger.Error("[Request.Get] %s", err.Error())
		return nil, err
	}
	t3 = time.Now().UnixNano()
	for _, cb := range callback {
		result.CallbackOutput = append(result.CallbackOutput, cb.Do(req, resp))
		cbCount++
	}
	result.status = resp.StatusCode()
	result.body = resp.Body()
	result.contentType = resp.Header.ContentType()
	t4 = time.Now().UnixNano()
	reqLogger.Infof("[Request.Get] target url: %s, transport time cost: %.3fms, total time cost: %.3fms", string(req.URI().FullURI()), float64(t3-t2)/1e6, float64(t4-t1)/1e6)
	reqLogger.Infof("[Request.Get] %d callback triggered, check response on HTTPResponse.CallbackOutput", cbCount)
	return &result, nil
}

func (r *requests) PostJson(url string, param interface{}, headers map[string]string, callback ...HTTPCallback) (*HTTPResponse, error) {
	var result HTTPResponse
	var t1, t2, t3, t4 int64
	var cbCount int

	t1 = time.Now().UnixNano()
	u, err := url2.Parse(url)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	reqUri := fasthttp.AcquireURI()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseURI(reqUri)

	reqUri.SetHost(u.Host)
	reqUri.SetPath(u.Path)
	reqUri.SetScheme(u.Scheme)
	req.SetRequestURIBytes(reqUri.FullURI())
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	body, err := json.Marshal(param)
	if err != nil {
		reqLogger.Error("[Request.Post] %s", err.Error())
		return nil, err
	}
	req.SetBody(body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	t2 = time.Now().UnixNano()
	err = fasthttp.Do(req, resp)
	if err != nil {
		reqLogger.Error("[Request.Post] %s", err.Error())
		return nil, err
	}
	t3 = time.Now().UnixNano()
	for _, cb := range callback {
		result.CallbackOutput = append(result.CallbackOutput, cb.Do(req, resp))
		cbCount++
	}
	result.status = resp.StatusCode()
	result.body = resp.Body()
	result.contentType = resp.Header.ContentType()
	t4 = time.Now().UnixNano()
	reqLogger.Infof("[Request.Post] target url: %s, transport time cost: %.3fms, total time cost: %.3fms", string(req.URI().FullURI()), float64(t3-t2)/1e6, float64(t4-t1)/1e6)
	reqLogger.Infof("[Request.Post] %d callback triggered, check response on HTTPResponse.CallbackOutput", cbCount)
	return &result, nil
}

func (r *requests) PostForm(url string, param map[string]string, headers map[string]string, callback ...HTTPCallback) (*HTTPResponse, error) {
	var result HTTPResponse
	var t1, t2, t3, t4 int64
	var cbCount int

	t1 = time.Now().UnixNano()
	u, err := url2.Parse(url)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	reqUri := fasthttp.AcquireURI()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseURI(reqUri)

	reqUri.SetHost(u.Host)
	reqUri.SetPath(u.Path)
	reqUri.SetScheme(u.Scheme)
	req.SetRequestURIBytes(reqUri.FullURI())
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	values := make(url2.Values)
	for k, v := range param {
		values[k] = []string{v}
	}
	body := values.Encode()
	req.SetBody([]byte(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	t2 = time.Now().UnixNano()
	err = fasthttp.Do(req, resp)
	if err != nil {
		reqLogger.Error("[Request.Post] %s", err.Error())
		return nil, err
	}
	t3 = time.Now().UnixNano()
	for _, cb := range callback {
		result.CallbackOutput = append(result.CallbackOutput, cb.Do(req, resp))
		cbCount++
	}
	result.status = resp.StatusCode()
	result.body = resp.Body()
	result.contentType = resp.Header.ContentType()
	t4 = time.Now().UnixNano()
	reqLogger.Infof("[Request.Post] target url: %s, transport time cost: %.3fms, total time cost: %.3fms", string(req.URI().FullURI()), float64(t3-t2)/1e6, float64(t4-t1)/1e6)
	reqLogger.Infof("[Request.Post] %d callback triggered, check response on HTTPResponse.CallbackOutput", cbCount)
	return &result, nil
}

func (r *requests) PostJsonWithQueryString(url string, param interface{}, headers map[string]string, query map[string]string, callback ...HTTPCallback) (*HTTPResponse, error) {
	var result HTTPResponse
	var t1, t2, t3, t4 int64
	var cbCount int

	t1 = time.Now().UnixNano()
	u, err := url2.Parse(url)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	req := fasthttp.AcquireRequest()
	reqUri := fasthttp.AcquireURI()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseURI(reqUri)

	reqUri.SetHost(u.Host)
	reqUri.SetPath(u.Path)
	reqUri.SetScheme(u.Scheme)
	queryValues := make(url2.Values)
	for k, v := range query {
		queryValues[k] = []string{v}
	}
	queryString := queryValues.Encode()
	reqUri.SetQueryString(queryString)
	req.SetRequestURIBytes(reqUri.FullURI())
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	body, err := json.Marshal(param)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	req.SetBody(body)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	t2 = time.Now().UnixNano()
	err = fasthttp.Do(req, resp)
	if err != nil {
		reqLogger.Exception(err)
		return nil, err
	}
	t3 = time.Now().UnixNano()
	for _, cb := range callback {
		result.CallbackOutput = append(result.CallbackOutput, cb.Do(req, resp))
		cbCount++
	}
	result.status = resp.StatusCode()
	result.body = resp.Body()
	result.contentType = resp.Header.ContentType()
	t4 = time.Now().UnixNano()
	reqLogger.Infof("[Request.Post] target url: %s, transport time cost: %.3fms, total time cost: %.3fms", string(req.URI().FullURI()), float64(t3-t2)/1e6, float64(t4-t1)/1e6)
	reqLogger.Infof("[Request.Post] %d callback triggered, check response on HTTPResponse.CallbackOutput", cbCount)
	return &result, nil
}
