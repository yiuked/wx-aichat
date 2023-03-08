package apis

import (
	"bios-dev/lib/wx"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AddDraft(w http.ResponseWriter, r *http.Request) {
	// 接收消息并回复
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error: ", err)
		return
	}
	var params wx.AddDraftReq
	log.Println("request=> ", string(body))
	err = xml.Unmarshal(body, &params)
	if err != nil {
		fmt.Println("xml.Unmarshal error: ", err)
		return
	}
	token, err := wx.GetAccessToken()
	if err != nil {
		return
	}

	wx.AddDraft(token.AccessToken, params)
}
