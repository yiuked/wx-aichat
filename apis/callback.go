// Package apis 微信回调
// 相关文档 https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Receiving_event_pushes.html
package apis

import (
	"bios-dev/apis/internet"
	"bios-dev/config"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	// 微信服务器验证
	signature := r.URL.Query().Get("signature")
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	if !checkSignature(config.WxToken, signature, timestamp, nonce) {
		fmt.Fprintf(w, "Invalid signature")
		return
	}

	if r.Method == "GET" {
		echostr := r.URL.Query().Get("echostr")
		fmt.Fprintf(w, echostr)
		return
	} else if r.Method == "POST" {
		// 接收消息并回复
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("ioutil.ReadAll error: ", err)
			return
		}
		var msg internet.Message
		log.Println("request=> ", string(body))
		err = xml.Unmarshal(body, &msg)
		if err != nil {
			fmt.Println("xml.Unmarshal error: ", err)
			return
		}
		ctx := &internet.MsgContext{Msg: &msg, Request: r, ResponseWriter: w}
		if msg.MsgType == "text" {
			internet.TextHandel(ctx)
		} else if msg.MsgType == "event" { // 如果是事件消息
			internet.EventHandel(ctx)
		}
	}
}

func checkSignature(token, signature, timestamp, nonce string) bool {
	strs := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(strs)

	sha1Str := sha1.New()
	_, _ = sha1Str.Write([]byte(strings.Join(strs, "")))
	encryptedStr := hex.EncodeToString(sha1Str.Sum(nil))

	// 加入时间检测
	atoi, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	if time.Now().Unix()-atoi > 60 {
		return false
	}

	return encryptedStr == signature
}
