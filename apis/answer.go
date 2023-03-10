package apis

import (
	"bios-dev/apis/internet"
	"net/http"
)

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		internet.GetAnswer(&internet.MsgContext{Request: r, ResponseWriter: w})
	}
}
