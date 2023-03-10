package apis

import (
	"bios-dev/apis/internet"
	"net/http"
)

func AddDraft(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		internet.AddDraft(&internet.MsgContext{Request: r, ResponseWriter: w})
	}
}
