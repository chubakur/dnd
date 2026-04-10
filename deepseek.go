package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	deepseek_base_url = "https://api.deepseek.com/chat/completions"
	deepseek_model    = "deepseek-chat"
)

type deepSeekClient struct {
	apiKey string
}

func NewDeepSeekClient(apiKey string) *deepSeekClient {
	return &deepSeekClient{apiKey: apiKey}
}

type deepSeekRoleContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type deepSeekQuery struct {
	Model    string                `json:"model"`
	Stream   bool                  `json:"stream"`
	Messages []deepSeekRoleContent `json:"messages"`
}

func NewLLMUserMessage(content string) deepSeekRoleContent {
	return deepSeekRoleContent{Role: "user", Content: content}
}

func NewLLMSystemMessage(content string) deepSeekRoleContent {
	return deepSeekRoleContent{Role: "system", Content: content}
}

func (c *deepSeekClient) Query(message string) (string, error) {
	body := deepSeekQuery{Model: deepseek_model, Stream: false, Messages: []deepSeekRoleContent{NewLLMUserMessage(message)}}
	json_body, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", deepseek_base_url, strings.NewReader(string(json_body)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(resBody), nil
}
