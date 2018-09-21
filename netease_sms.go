package utils

import (
	"bytes"
	"crypto/sha1"
	"github.com/riposa/utils/errors"
	"math/rand"
	"strconv"
	"time"
)

const (

	InterfaceSendSmsCode   = "https://api.netease.im/sms/sendcode.action"
	InterfaceVerifySmsCode = "https://api.netease.im/sms/verifycode.action"
)

func GetRandomString(l int) string {
	var result []byte
	seedByte := []byte(seed)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, seedByte[r.Intn(len(seedByte))])
	}
	return string(result)
}

func ByteToHex(data [20]byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}
	return buffer.String()
}

func makeCheckSum(nonce, curTime string) string {
	hashed := sha1.Sum([]byte(NeteaseAppSecret + nonce + curTime))
	return ByteToHex(hashed)
}

func requestHeader() map[string]string {
	nonce := GetRandomString(50)
	curTime := strconv.FormatUint(uint64(time.Now().Unix()), 10)
	return map[string]string{
		"AppKey":   NeteaseAppKey,
		"Nonce":    nonce,
		"CurTime":  curTime,
		"CheckSum": makeCheckSum(nonce, curTime),
	}
}

func SendSmsCode(phone string) (bool, error) {
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Obj  string `json:"obj"`
	}
	header := requestHeader()
	header["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := Requests.PostForm(
		InterfaceSendSmsCode,
		map[string]string{
			"mobile":     phone,
			"templateid": SmsTemplateID,
		},
		header,
	)
	if err != nil {
		return false, err
	}
	resp.Json(&res)
	if res.Code != 200 {
		return false, errors.NewFormat(80, res.Msg)
	}
	return true, nil
}

func VerifySmsCode(phone, code string) (bool, error) {
	var res struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Obj  string `json:"obj"`
	}
	header := requestHeader()
	header["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := Requests.PostForm(
		InterfaceVerifySmsCode,
		map[string]string{
			"mobile": phone,
			"code":   code,
		},
		header,
	)
	if err != nil {
		return false, err
	}
	resp.Json(&res)
	if res.Code != 200 {
		return false, errors.NewFormat(80, res.Msg)
	}
	return true, nil
}
