package llmcore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chubakur/dnd/mcp"
	"github.com/chubakur/dnd/types"
)

const (
	deepseek_base_url = "https://api.deepseek.com/chat/completions"
	deepseek_model    = "deepseek-chat"
)

type deepSeekClient struct {
	apiKey string
	tools  []*mcp.MCPTool
}

func NewDeepSeekClient(apiKey string, mcpTools []*mcp.MCPTool) *deepSeekClient {
	return &deepSeekClient{apiKey: apiKey, tools: mcpTools}
}

type DeepSeekRoleContent struct {
	Role       string                           `json:"role"`
	Content    string                           `json:"content"`
	ToolCalls  []types.DeepseekResponseToolCall `json:"tool_calls,omitempty"`
	ToolCallId string                           `json:"tool_call_id,omitempty"`
}

type deepSeekQuery struct {
	Model    string                 `json:"model"`
	Stream   bool                   `json:"stream"`
	Messages []DeepSeekRoleContent  `json:"messages"`
	Tools    []types.WrappedMCPTool `json:"tools"`
}

type deepseekResponseChoice struct {
	Index        int                 `json:"index"`
	Message      DeepSeekRoleContent `json:"message"`
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

func wrapMcp(mc *mcp.MCPTool) types.WrappedMCPTool {
	result := types.WrappedMCPTool{
		Type: "function",
		Function: types.WrappedMCPFunction{
			Type:        mc.Parameters.Type,
			Name:        mc.Name,
			Description: mc.Description,
			Parameters: types.WrappedMCPFunctionParameters{
				Type:       "object",
				Properties: make(map[string]types.WrappedMCPFunctionParametersProperty),
				Required:   make([]string, 0),
			},
		},
	}
	for _, param := range mc.Parameters.Properties {
		result.Function.Parameters.Properties[param.Name] = types.WrappedMCPFunctionParametersProperty{
			Type:        param.Type,
			Description: param.Description,
		}
		if param.IsRequired {
			result.Function.Parameters.Required = append(result.Function.Parameters.Required, param.Name)
		}
	}
	return result
}

func NewLLMUserMessage(content string) DeepSeekRoleContent {
	return DeepSeekRoleContent{Role: "user", Content: content}
}

func (c *deepSeekClient) Query(mc *MessageChain) (*deepseekResponse, error) {
	wrappedTools := make([]types.WrappedMCPTool, len(c.tools))
	for i, tool := range c.tools {
		wrappedTools[i] = wrapMcp(tool)
	}
	body := deepSeekQuery{
		Model:    deepseek_model,
		Stream:   false,
		Messages: mc.chain,
		Tools:    wrappedTools,
	}
	json_body, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", deepseek_base_url, strings.NewReader(string(json_body)))
	fmt.Println(string(json_body))
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
	deepseekResult := deepseekResponse{}
	// TODO: decode errors if not 200
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&deepseekResult)
	if err != nil {
		return nil, err
	}
	return &deepseekResult, nil
}
