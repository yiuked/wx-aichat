package apis

import (
	"log"
	"net/http"
	"wx-aichat/lib/wx"
)

func Publish(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token, _ := wx.GetAccessToken()
		log.Println(token)
	}
}
