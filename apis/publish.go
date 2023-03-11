package apis

import (
	"bios-dev/lib/wx"
	"log"
	"net/http"
)

func Publish(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		token, _ := wx.GetAccessToken()
		log.Println(token)
	}
}
