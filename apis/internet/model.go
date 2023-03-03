package internet

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	EventSubscribe   string = "subscribe"   // 定阅
	EventUnsubscribe string = "unsubscribe" // 取消定阅
)

const (
	MsgText  string = "text"
	MsgImage        = "image"
	MsgVoice        = "voice"
	MsgMusic        = "music"
	MsgNews         = "news"
)

type MsgContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Msg            *Message
	// 用于存储键值对的 map
	data map[string]interface{}
	// parent 用于实现 context 树
	parent context.Context
}

type Message struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Event        string
	EventKey     string
	Content      string
	MsgId        int64
}

func (m Message) Collection() string {
	return "message"
}

type Response struct {
	XMLName      struct{} `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}

type Faq struct {
	UserName   string
	CreateTime int64
	Question   string
	Answer     string
}

func (mc *MsgContext) Value(key interface{}) interface{} {
	// 先查找本地的键值对
	val, ok := mc.data[key.(string)]
	if ok {
		return val
	}
	// 如果本地找不到，就向上查找
	return mc.parent.Value(key)
}

func (mc *MsgContext) WithValue(key interface{}, val interface{}) context.Context {
	// 创建一个新的 map
	data := make(map[string]interface{})
	// 将新的键值对存入 map
	data[key.(string)] = val
	// 创建一个新的 MyContext
	return &MsgContext{
		data:   data,
		parent: mc,
	}
}

func (mc *MsgContext) Deadline() (deadline time.Time, ok bool) {
	// 由于该示例不需要实现超时控制，所以这里直接返回 false
	return time.Time{}, false
}

func (mc *MsgContext) Done() <-chan struct{} {
	// 由于该示例不需要实现超时控制，所以这里直接返回 nil
	return nil
}

func (mc *MsgContext) Err() error {
	// 由于该示例不需要实现超时控制，所以这里直接返回 nil
	return nil
}

func Error(w io.Writer) {
	fmt.Fprintf(w, "系统繁忙，请稍后重试！")
}

func TodayLimit(w io.Writer) {
	fmt.Fprintf(w, "今日已超限！")
}
