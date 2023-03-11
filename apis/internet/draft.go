package internet

import (
	"bios-dev/config"
	"bios-dev/lib/wx"
	"context"
	"github.com/russross/blackfriday/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddDraft(ctx *MsgContext) {
	// 接收消息并回复
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
	var faq Faq
	if err := config.Mgo.Db.Collection("faq").FindOne(context.Background(), bson.M{"_id": hex}).
		Decode(&faq); err == nil {
		if faq.UserName == config.WXOpenId {
			var params wx.AddDraftReq
			html := blackfriday.Run([]byte(faq.Answer))
			params.Articles = append(params.Articles, wx.Article{
				ThumbMediaID:       config.WxMediaId,
				Author:             config.WxAuthor,
				OnlyFansCanComment: 0,
				NeedOpenComment:    0,
				Title:              faq.Question,
				Content:            string(html),
			})

			token, err := wx.GetAccessToken()
			if err != nil {
				ResponseText(ctx.ResponseWriter, "获取微信token失败")
				return
			}

			wx.AddDraft(token.AccessToken, params)
			ResponseText(ctx.ResponseWriter, "发布成功")
		} else {
			ResponseText(ctx.ResponseWriter, "权限不足")
		}
		return
	}
}
