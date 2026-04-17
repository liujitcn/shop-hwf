package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type WxSessionKey struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}
type WxAccessToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
type PhoneNumber struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int    `json:"timestamp"`
			Appid     string `json:"appid"`
		} `json:"watermark"`
	} `json:"phone_info"`
}

func GetPhoneNumber(accessToken, code string) (*PhoneNumber, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req := make(map[string]string)
	req["code"] = code
	reqB, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	resp, err = client.Post(url, "application/json", strings.NewReader(string(reqB)))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var b PhoneNumber
	err = json.Unmarshal(body, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func GetAccessToken(appid, secret string) (*WxAccessToken, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var b WxAccessToken
	err = json.Unmarshal(body, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func GetSessionKey(appid, secret, code string) (*WxSessionKey, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var b WxSessionKey
	err = json.Unmarshal(body, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
