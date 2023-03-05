package internet

import (
	"bios-dev/config"
	"bios-dev/lib"
	"context"
	"encoding/xml"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strings"
	"sync"
	"time"
)

var locks map[string]*sync.Mutex

func init() {
	locks = make(map[string]*sync.Mutex)
}

func TextHandel(ctx *MsgContext) {
	if v, b := config.Limit[ctx.Msg.FromUserName]; b {
		if v.Cnt > 30 {
			TodayLimit(ctx.ResponseWriter)
			return
		}
		if time.Now().Unix()-v.LastSt <= 3 {
			Error(ctx.ResponseWriter)
			return
		}
	} else {
		config.Limit[ctx.Msg.FromUserName] = &config.UserLimit{Cnt: 0, LastSt: time.Now().Unix()}
	}

	if _, ok := locks[ctx.Msg.FromUserName]; !ok {
		locks[ctx.Msg.FromUserName] = &sync.Mutex{}

	}
	// 加锁，保障只有一次请求到openai
	locks[ctx.Msg.FromUserName].Lock()
	defer locks[ctx.Msg.FromUserName].Unlock()

	// 查询结果是否存在，存在则直接返回
	var existsFaq Faq
	if err := config.Mgo.Db.Collection("faq").FindOne(context.Background(), bson.M{"question": ctx.Msg.Content}).
		Decode(&existsFaq); err == nil {
		resp := Response{
			ToUserName:   ctx.Msg.FromUserName,
			FromUserName: ctx.Msg.ToUserName,
			CreateTime:   ctx.Msg.CreateTime,
			MsgType:      "text",
			Content:      existsFaq.Answer,
		}
		res, err := xml.Marshal(resp)
		if err != nil {
			fmt.Println("xml.Marshal error: ", err)
			Error(ctx.ResponseWriter)
			return
		}
		log.Println("response=> ", string(res))

		config.Limit[ctx.Msg.FromUserName].Cnt += 1
		config.Limit[ctx.Msg.FromUserName].LastSt = time.Now().Unix()
		fmt.Fprintf(ctx.ResponseWriter, string(res))
		return
	}

	result := lib.Send(ctx.Msg.Content)
	if result == nil {
		Error(ctx.ResponseWriter)
	}

	var returnMsg string
	for _, msg := range result {
		for _, m := range msg.Choices {
			returnMsg = returnMsg + strings.Trim(m.Message.Content, "\n")
		}
	}

	resp := Response{
		ToUserName:   ctx.Msg.FromUserName,
		FromUserName: ctx.Msg.ToUserName,
		CreateTime:   ctx.Msg.CreateTime,
		MsgType:      "text",
		Content:      returnMsg,
	}

	faq := Faq{UserName: ctx.Msg.FromUserName, CreateTime: time.Now().Unix(), Question: ctx.Msg.Content, Answer: resp.Content}
	if _, err := config.Mgo.Db.Collection("faq").InsertOne(context.Background(), &faq); err != nil {
		log.Println("mongodb insert err,", err)
		Error(ctx.ResponseWriter)
		return
	}

	res, err := xml.Marshal(resp)
	if err != nil {
		fmt.Println("xml.Marshal error: ", err)
		Error(ctx.ResponseWriter)
		return
	}
	log.Println("response=> ", string(res))

	config.Limit[ctx.Msg.FromUserName].Cnt += 1
	config.Limit[ctx.Msg.FromUserName].LastSt = time.Now().Unix()
	fmt.Fprintf(ctx.ResponseWriter, string(res))
}
