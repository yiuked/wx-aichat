package main

import (
	"bios-dev/apis"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var jsFiles embed.FS

func main() {
	http.HandleFunc("/callback", apis.Callback)
	http.HandleFunc("/getAnswer", apis.GetAnswer)
	http.HandleFunc("/addDraft", apis.AddDraft)
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
	log.Println("HTTP server run success http://localhost:8089")
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe error: ", err)
	}
}
