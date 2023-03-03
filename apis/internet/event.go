package internet

import (
	"encoding/xml"
	"log"
	"time"
)

func EventHandel(ctx *MsgContext) {
	if ctx.Msg.Event == "subscribe" { // 如果是关注事件
		// 构造回复的消息
		replyMsg := Response{
			ToUserName:   ctx.Msg.FromUserName,
			FromUserName: ctx.Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "一道灵光咋现，你发现了一个神奇洞穴，它似乎拥有神奇的力量，知晓所有的秘密和答案，它每天可以回答30个问题。",
		}
		// 将回复的消息序列化为XML格式
		replyMsgXml, err := xml.Marshal(replyMsg)
		if err != nil {
			log.Println(err)
			Error(ctx.ResponseWriter)
			return
		}
		// 将XML格式的回复消息返回给微信服务器
		ctx.ResponseWriter.Write(replyMsgXml)
	} else if ctx.Msg.Event == "unsubscribe" { // 如果是取消关注事件
		// 构造回复的消息
		replyMsg := Response{
			ToUserName:   ctx.Msg.FromUserName,
			FromUserName: ctx.Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "一道青烟冒出，又一位道友飞升离开。",
		}
		// 将回复的消息序列化为XML格式
		replyMsgXml, err := xml.Marshal(replyMsg)
		if err != nil {
			log.Println(err)
			Error(ctx.ResponseWriter)
			return
		}
		// 将XML格式的回复消息返回给微信服务器
		ctx.ResponseWriter.Write(replyMsgXml)
	}
}
