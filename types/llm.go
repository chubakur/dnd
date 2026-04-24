package types

type DeepSeekRoleContent struct {
	Role       string                           `json:"role"`
	Content    string                           `json:"content"`
	ToolCalls  []DeepseekResponseToolCall `json:"tool_calls,omitempty"`
	ToolCallId string                           `json:"tool_call_id,omitempty"`
}
