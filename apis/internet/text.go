package internet

import (
	"bios-dev/config"
	"bios-dev/lib"
	"bios-dev/lib/wx"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"strings"
	"time"
)

var locks map[string]int

func init() {
	locks = make(map[string]int)
}

func TextHandel(ctx *MsgContext) {
	resp := Response{
		ToUserName:   ctx.Msg.FromUserName,
		FromUserName: ctx.Msg.ToUserName,
		CreateTime:   ctx.Msg.CreateTime,
		MsgType:      "text",
	}
	// 检测是否有超过当日30次限制
	if v, b := config.Limit[ctx.Msg.FromUserName]; b {
		if v.Cnt > 30 {
			TodayLimit(ctx.ResponseWriter, resp)
			return
		}
	} else {
		config.Limit[ctx.Msg.FromUserName] = &config.UserLimit{Cnt: 0, LastSt: time.Now().Unix()}
	}

	taskKey := lib.MD5(fmt.Sprintf("%s:%s", ctx.Msg.FromUserName, ctx.Msg.Content))

	if _, ok := locks[taskKey]; !ok {
		locks[taskKey] = 0
	}
	locks[taskKey] += 1

	// 查询结果是否存在，存在则直接返回
	var existsFaq Faq
	if err := config.Mgo.Db.Collection("faq").FindOne(context.Background(), bson.M{"question": ctx.Msg.Content}).
		Decode(&existsFaq); err == nil {
		if len(existsFaq.Answer) > 1000 {
			resp.Content = fmt.Sprintf(`<a href="%s/getAnswer?uuid=%s">啊哈！篇幅过长拉，点击查看详情吧</a>`, config.HOST, existsFaq.Uuid.Hex())
		} else {
			resp.Content = existsFaq.Answer
		}
		delete(locks, taskKey)
		ResponseXML(ctx.ResponseWriter, resp)
		return
	}
	// 第一次请求过来，转到openai去查，如果5秒内有回复，直接返给微信，如果超过5秒，微信发起第二次请求
	// 第二次请求过来，如果第一次已经有结果并写到数据库了，那么上面的的代码执行，结果直接返回，否则延时5秒，微信会发起第三次请求。
	// 第三次请求过来，如果第一次已经有结果并写到数据库了，那么上面的的代码执行，结果直接返回，否则累计已经过了10秒，无需再等直接返回超时重试。
	if locks[taskKey] == 2 {
		time.Sleep(5 * time.Second)
		return
	} else if locks[taskKey] > 2 {
		fmt.Println("请求超时")
		resp.Content = fmt.Sprintf(`<a href="weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=1">糟糕！断片了，点我重试一次</a>`, ctx.Msg.Content)
		if locks[taskKey] >= 3 {
			delete(locks, taskKey)
		}
		ResponseXML(ctx.ResponseWriter, resp)
		return
	}
	// 只有第一次会执行以下代码
	result := wx.Send(ctx.Msg.Content)
	log.Println("收到结果")
	var returnMsg string
	for _, msg := range result {
		for _, m := range msg.Choices {
			returnMsg = returnMsg + strings.Trim(m.Message.Content, "\n")
		}
	}

	faq := Faq{
		Uuid:       primitive.NewObjectID(),
		UserName:   ctx.Msg.FromUserName,
		CreateTime: time.Now().Unix(),
		Question:   ctx.Msg.Content,
		Answer:     returnMsg,
	}
	if _, err := config.Mgo.Db.Collection("faq").InsertOne(context.Background(), &faq); err != nil {
		log.Println("mongodb insert err,", err)
		Error(ctx.ResponseWriter, resp)
		return
	}

	if len(faq.Answer) > 1000 {
		resp.Content = fmt.Sprintf(`<a href="%s/getAnswer?uuid=%s">啊哈！篇幅过长拉，点击查看详情吧</a>`, config.HOST, faq.Uuid.Hex())
	} else {
		resp.Content = faq.Answer
	}

	config.Limit[ctx.Msg.FromUserName].Cnt += 1
	config.Limit[ctx.Msg.FromUserName].LastSt = time.Now().Unix()
	delete(locks, taskKey)
	ResponseXML(ctx.ResponseWriter, resp)
}
