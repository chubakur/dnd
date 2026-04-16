package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	deepseek_base_url = "https://api.deepseek.com/chat/completions"
	deepseek_model    = "deepseek-chat"
)

type deepSeekClient struct {
	apiKey string
	tools  []*MCPTool
}

func NewDeepSeekClient(apiKey string, mcpTools []*MCPTool) *deepSeekClient {
	return &deepSeekClient{apiKey: apiKey, tools: mcpTools}
}

type deepseekResponseToolCall struct {
	Index    int    `json:"index"`
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type deepSeekRoleContent struct {
	Role      string                     `json:"role"`
	Content   string                     `json:"content"`
	ToolCalls []deepseekResponseToolCall `json:"tool_calls,omitempty"`
}

type deepSeekQuery struct {
	Model    string                `json:"model"`
	Stream   bool                  `json:"stream"`
	Messages []deepSeekRoleContent `json:"messages"`
	Tools    []WrappedMCPTool      `json:"tools"`
}

type deepseekResponseChoice struct {
	Index        int                 `json:"index"`
	Message      deepSeekRoleContent `json:"message"`
	Logprobs     any                 `json:"logprobs"`
	FinishReason string              `json:"finish_reason"`
}

type deepseekUsageStat struct {
	PromptTokens        int `json:"prompt_tokens"`
	CompletionTokens    int `json:"complection_tokens"`
	TotalTokens         int `json:"total_tokens"`
	PromptTokensDetails struct {
		CachedTokens int `json:"cached_tokens"`
	} `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}

type deepseekResponse struct {
	Id                string                   `json:"id"`
	Object            string                   `json:"object"`
	Created           int                      `json:"created"`
	Model             string                   `json:"model"`
	Choices           []deepseekResponseChoice `json:"choices"`
	Usage             deepseekUsageStat        `json:"usage"`
	SystemFingerprint string                   `json:"system_fingerprint"`
}

func NewLLMUserMessage(content string) deepSeekRoleContent {
	return deepSeekRoleContent{Role: "user", Content: content}
}

func NewLLMSystemMessage(content string) deepSeekRoleContent {
	return deepSeekRoleContent{Role: "system", Content: content}
}

func (c *deepSeekClient) Query(message string) (*deepseekResponse, error) {
	wrappedTools := make([]WrappedMCPTool, len(c.tools))
	for i, tool := range c.tools {
		wrappedTools[i] = WrappedMCPTool{Type: "function", Function: tool}
	}
	body := deepSeekQuery{
		Model:    deepseek_model,
		Stream:   false,
		Messages: []deepSeekRoleContent{NewLLMUserMessage(message)},
		Tools:    wrappedTools,
	}
	json_body, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", deepseek_base_url, strings.NewReader(string(json_body)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// resBody, err := io.ReadAll(resp.Body)
	deepseekResult := deepseekResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&deepseekResult)
	if err != nil {
		return nil, err
	}
	return &deepseekResult, nil
}
