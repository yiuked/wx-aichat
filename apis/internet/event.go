package internet

import (
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
			Content:      "一道灵光咋现，你发现了一个神奇洞穴，它似乎拥有神奇的力量，知晓所有的秘密和答案，它每天可以回答50个问题。",
		}
		ResponseXML(ctx.ResponseWriter, replyMsg)
	} else if ctx.Msg.Event == "unsubscribe" { // 如果是取消关注事件
		// 构造回复的消息
		replyMsg := Response{
			ToUserName:   ctx.Msg.FromUserName,
			FromUserName: ctx.Msg.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      "text",
			Content:      "一道青烟冒出，又一位道友飞升离开。",
		}
		ResponseXML(ctx.ResponseWriter, replyMsg)
	}
}
