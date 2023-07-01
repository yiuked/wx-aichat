package internet

import (
	"context"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
	"sync"
	"time"
	"wx-aichat/config"
	"wx-aichat/lib/wx"
)

var addDraftLock sync.Mutex

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

	cacheKey := fmt.Sprintf("draft%s", uuid)
	addDraftLock.Lock()
	defer addDraftLock.Unlock()

	if _, err := config.Cache.Get(cacheKey); err == nil {
		ResponseText(ctx.ResponseWriter, "操作频繁")
		return
	}
	// 微信代码块样式
	// <pre class="code-snippet" data-lang="python"></pre>
	var faq Faq
	if err := config.Mgo.Db.Collection("faq").FindOne(context.Background(), bson.M{"_id": hex}).
		Decode(&faq); err == nil {
		if faq.UserName == config.WXOpenId {
			var params wx.AddDraftReq
			html := blackfriday.Run([]byte(faq.Answer))
			re := regexp.MustCompile(`<pre><code class="language-(.*?)">`)
			htmlString := string(html)
			htmlString = re.ReplaceAllString(htmlString, `<pre class="code-snippet code-snippet_nowrap code-snippet__$1" data-lang="$1"><code>`)
			params.Articles = append(params.Articles, wx.Article{
				ThumbMediaID:       config.WxMediaId,
				Author:             config.WxAuthor,
				OnlyFansCanComment: 0,
				NeedOpenComment:    0,
				Title:              faq.Question,
				Content:            htmlString,
			})

			token, err := wx.GetAccessToken()
			if err != nil {
				ResponseText(ctx.ResponseWriter, "获取微信token失败")
				return
			}

			wx.AddDraft(token.AccessToken, params)
			config.Cache.Set(cacheKey, 1, 60*time.Second)
			ResponseText(ctx.ResponseWriter, "发布成功")
		} else {
			ResponseText(ctx.ResponseWriter, "权限不足")
		}
		return
	}
}
