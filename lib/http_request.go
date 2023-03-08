package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func Request(method, url string, header http.Header, data interface{}) (result []byte, err error) {
	RequestLog("URL：%s", url)
	var body io.Reader
	var marshal []byte
	if data != nil {
		marshal, err = json.Marshal(data)
		if err != nil {
			RequestLog("Fail json marshal err：%v", err)
			return
		}
		body = bytes.NewReader(marshal)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		RequestLog("Fail http new request err：%v", err)
		return
	}
	req.Header = header

	// 打印日志
	RequestLog("Header")
	for key, val := range req.Header {
		RequestLog("%s: %s", key, strings.Join(val, ","))
	}
	RequestLog("\n%s\n", string(marshal))

	// 发送 HTTP 请求
	client := &http.Client{}
	var tryCnt int
TRY:
	res, err := client.Do(req)
	if err != nil {
		// 重试3次
		if tryCnt < 3 {
			tryCnt++
			RequestLog("Failed to send HTTP request:", err)
			RequestLog("Retry %d", tryCnt)
			goto TRY
		}
		RequestLog("Retry %d fail exit", tryCnt)
		return
	}

	// 解析 HTTP 响应
	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		RequestLog("Failed to read response body:%v", err)
		return
	}

	log.Println(string(all))
	return all, nil
}

func RequestLog(format string, v ...any) {
	log.Printf(format+"\n", v...)
}
