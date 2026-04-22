package types

type MCPTool struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  MCPToolParameters   `json:"parameters"`
	F           MCPToolFunc
}

type MCPToolFunc func(*Transport, *DeepseekResponseToolCall) (string, error)

type MCPToolParameters struct {
	Type       string
	Properties []MCPToolProperty
}

type MCPToolProperty struct {
	Name        string
	Type        string
	Description string
	IsRequired  bool
}

type DeepseekResponseToolCall struct {
	Index    int    `json:"index"`
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type WrappedMCPTool struct {
	Type     string                  `json:"type"`
	Function WrappedMCPFunction      `json:"function"`
}

type WrappedMCPFunction struct {
	Type        string                       `json:"type"`
	Name        string                       `json:"name"`
	Description string                       `json:"description"`
	Parameters  WrappedMCPFunctionParameters `json:"parameters"`
}

type WrappedMCPFunctionParametersProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type WrappedMCPFunctionParameters struct {
	Type       string                                    `json:"type"`
	Required   []string                                `json:"required"`
	Properties map[string]WrappedMCPFunctionParametersProperty `json:"properties"`
}

type MCPResult struct {
	Function   string
	Result     string
	Error      error
	ToolCallId string
}

type Transport struct {
	YdbClient interface{}
	Ctx       interface{}
}

func WrapMCPTool(mc *MCPTool) WrappedMCPTool {
	result := WrappedMCPTool{
		Type: "function",
		Function: WrappedMCPFunction{
			Type:        mc.Parameters.Type,
			Name:        mc.Name,
			Description: mc.Description,
			Parameters: WrappedMCPFunctionParameters{
				Type:       "object",
				Properties: make(map[string]WrappedMCPFunctionParametersProperty),
				Required:   make([]string, 0),
			},
		},
	}

	for _, param := range mc.Parameters.Properties {
		result.Function.Parameters.Properties[param.Name] = WrappedMCPFunctionParametersProperty{
			Type:        param.Type,
			Description: param.Description,
		}
		if param.IsRequired {
			result.Function.Parameters.Required = append(result.Function.Parameters.Required, param.Name)
		}
	}

	return result
}
