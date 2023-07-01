package apis

import (
	"net/http"
	"wx-aichat/apis/internet"
)

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		internet.GetAnswer(&internet.MsgContext{Request: r, ResponseWriter: w})
	}
}
