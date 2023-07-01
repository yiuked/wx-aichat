package apis

import (
	"net/http"
	"wx-aichat/apis/internet"
)

func AddDraft(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		internet.AddDraft(&internet.MsgContext{Request: r, ResponseWriter: w})
	}
}
