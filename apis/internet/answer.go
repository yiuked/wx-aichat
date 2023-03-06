package internet

import (
	"bios-dev/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAnswer(ctx *MsgContext) {
	uuid := ctx.Request.URL.Query().Get("uuid")
	if len(uuid) <= 0 {
		ResponseText(ctx.ResponseWriter, "参数错误")
		return
	}
	hex, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		ResponseText(ctx.ResponseWriter, "参数错误")
		return
	}
	var existsFaq Faq
	if err := config.Mgo.Db.Collection("faq").FindOne(context.Background(), bson.M{"_id": hex}).
		Decode(&existsFaq); err == nil {
		ResponseText(ctx.ResponseWriter, `<script src="/md-page.js"></script><noscript>`+existsFaq.Answer)
		return
	}
	ResponseText(ctx.ResponseWriter, "读取失败")
	return
}
