package main

import (
	"bios-dev/apis"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/callback", apis.Callback)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe error: ", err)
	}
}
