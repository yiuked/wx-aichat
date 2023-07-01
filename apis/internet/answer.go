package internet

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"wx-aichat/config"
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
		if existsFaq.UserName == config.WXOpenId {
			if len(config.WxSecret) > 0 &&
				len(config.WxAppid) > 0 &&
				len(config.WXOpenId) > 0 &&
				len(config.WxMediaId) > 0 {
				ResponseText(ctx.ResponseWriter,
					fmt.Sprintf(`<script src="/md-page.js"></script><noscript>%s<a href="/addDraft?uuid=%s">发到草稿</a>`,
						existsFaq.Answer, existsFaq.Uuid.Hex()),
				)
			} else {
				ResponseText(ctx.ResponseWriter,
					fmt.Sprintf(`<script src="/md-page.js"></script><noscript>%s`, existsFaq.Answer),
				)
			}

		} else {
			ResponseText(ctx.ResponseWriter, fmt.Sprintf(`<script src="/md-page.js"></script><noscript>%s`, existsFaq.Answer))
		}
		return
	}
	ResponseText(ctx.ResponseWriter, "读取失败")
	return
}
