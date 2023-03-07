package apis

import (
	"bios-dev/apis/internet"
	"bios-dev/lib/wx"
	"net/http"
)

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		wx.GetAccessToken()
		internet.GetAnswer(&internet.MsgContext{Request: r, ResponseWriter: w})
	}
}
