package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wx-aichat/config"
)

const (
	RoleSystem    string = "system"
	RoleUser             = "user"
	RoleAssistant        = "assistant"
)

type Result struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClassificationRequest struct {
	Model   string    `json:"model"`
	Message []Message `json:"messages"`
}

type PaginationResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

func Send(msg string) []Result {
	// 构建 OpenAI API 请求 URL 和请求参数
	url := "https://api.openai.com/v1/chat/completions"
	reqBody, _ := json.Marshal(&ClassificationRequest{
		Model: "gpt-3.5-turbo",
		Message: []Message{
			{Role: RoleUser, Content: msg},
		},
	})

	// 构建 HTTP 请求对象
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.OpenAIKey))

	// 发送 HTTP 请求
	client := &http.Client{}
	var allResults []Result
	for {
		// 发送 HTTP 请求
		res, err := client.Do(req)
		if err != nil {
			log.Println("Failed to send HTTP request:", err)
			return nil
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("HTTP request fail %+v", res)
			return nil
		}

		// 解析 HTTP 响应
		defer res.Body.Close()
		var result Result
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			fmt.Println("Failed to decode HTTP response:", err)
			return nil
		}

		// 将当前请求返回的结果加入到全部结果中
		allResults = append(allResults, result)

		// 检查是否还有下一页，如果没有，退出循环
		pagination := res.Header.Get("Pagination")
		if len(pagination) > 0 {
			var paginationResponse PaginationResponse
			json.Unmarshal([]byte(pagination), &paginationResponse)
			if paginationResponse.Next == "" {
				break
			}

			// 更新下一个请求的 URL
			req.URL.RawQuery = paginationResponse.Next
		} else {
			break
		}
	}

	return allResults
}
