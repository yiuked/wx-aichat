package wx

import (
	"bios-dev/lib"
	"fmt"
	"net/http"
)

const (
	AddDraftApi = " https://api.weixin.qq.com/cgi-bin/draft/add?access_token=%s"
)

type AddDraftReq struct {
	Articles []struct {
		ThumbMediaID       string `json:"thumb_media_id"`
		Author             string `json:"author"`
		OnlyFansCanComment int    `json:"only_fans_can_comment"`
		Digest             string `json:"digest"`
		ContentSourceUrl   string `json:"content_source_url"`
		NeedOpenComment    int    `json:"need_open_comment"`
		Title              string `json:"title"`
		Content            string `json:"content"`
	} `json:"articles"`
}

// AddDraft 新建草稿
// https://developers.weixin.qq.com/doc/offiaccount/Draft_Box/Add_draft.html
func AddDraft(token string, data AddDraftReq) {
	headers := http.Header{}
	headers.Set("Content-type", "application/json")
	lib.Request(http.MethodPost, fmt.Sprintf(AddDraftApi, token), headers, data)
}
