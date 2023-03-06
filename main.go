package main

import (
	"bios-dev/apis"
	"bios-dev/config"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed static/*
var jsFiles embed.FS

func main() {
	// 每天24点清空 map 内所有值
	ticker := config.Tick24()

	http.HandleFunc("/callback", apis.Callback)
	http.HandleFunc("/getAnswer", apis.GetAnswer)
	// 将嵌入的 js/test.js 文件作为 HTTP 响应
	http.HandleFunc("/md-page.js", func(w http.ResponseWriter, r *http.Request) {
		content, err := fs.ReadFile(jsFiles, "static/md-page.js")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(content)
	})
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe error: ", err)
	}

	go func() {
		for {
			select {
			case <-ticker:
				config.ClearMap()
			}
		}
	}()
}
