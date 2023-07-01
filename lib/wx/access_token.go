package wx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"wx-aichat/config"
)

type AccessToken struct {
	ErrCode     int    `json:"errcode,omitempty"`
	ErrMsg      string `json:"errmsg,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

var accessToken = "access:token"

// GetAccessToken 获取ACCESS_TOKEN
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func GetAccessToken() (*AccessToken, error) {
	if get, err := config.Cache.Get(accessToken); err == nil {
		return get.(*AccessToken), nil
	}

	// 构建 HTTP 请求对象
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", config.WxAppid, config.WxSecret)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send HTTP request:", err)
		return nil, err
	}

	// 解析 HTTP 响应
	defer res.Body.Close()
	var result AccessToken
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Println("Failed to JSON decode:", err)
		return nil, err
	}
	if result.ErrCode > 0 {
		log.Println("Failed to get access_token:", result.ErrMsg)
		return nil, fmt.Errorf("code:%d,msg:%s", result.ErrCode, result.ErrMsg)
	}
	config.Cache.Set(accessToken, &result, 7000*time.Second)
	return &result, nil
}
