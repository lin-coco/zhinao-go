package zhinao

// Message 聊天消息
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // 消息内容
	Name       string     `json:"name,omitempty"`         // 消息发送者名称（可选）
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 工具调用列表（仅助手消息）
	ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID（仅工具消息）
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model             string      `json:"model"`                        // 模型名称（必需）
	Messages          []Message   `json:"messages"`                     // 消息列表（必需）
	Stream            bool        `json:"stream,omitempty"`             // 是否流式输出
	Temperature       float64     `json:"temperature,omitempty"`        // 温度参数 [0, 1]，默认0.9
	MaxTokens         int         `json:"max_tokens,omitempty"`         // 最大生成token数 [1, 2048]，默认2048
	TopP              float64     `json:"top_p,omitempty"`              // 核采样参数 [0, 1]，默认0.5
	TopK              int         `json:"top_k,omitempty"`              // TopK参数 [0, 1024]，默认0
	RepetitionPenalty float64     `json:"repetition_penalty,omitempty"` // 重复惩罚 [1, 2]，默认1.05
	NumBeams          int         `json:"num_beams,omitempty"`          // beam search数量 [1, 5]，默认1
	Tools             []Tool      `json:"tools,omitempty"`              // 工具列表
	ToolChoice        interface{} `json:"tool_choice,omitempty"`        // 工具选择，可以是字符串或结构体
	User              string      `json:"user,omitempty"`               // 用户标识
}

// Validate 验证请求参数
func (r *ChatRequest) Validate() error {
	if r.Model == "" {
		return ErrInvalidModel
	}
	if len(r.Messages) == 0 {
		return ErrEmptyMessages
	}
	if r.Temperature < 0 || r.Temperature > 1 {
		return &ValidationError{Field: "temperature", Message: "must be between 0 and 1"}
	}
	if r.TopP < 0 || r.TopP > 1 {
		return &ValidationError{Field: "top_p", Message: "must be between 0 and 1"}
	}
	return nil
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 响应选择项
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// ChatStreamResponse 流式聊天响应
type ChatStreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice 流式响应选择项
type StreamChoice struct {
	Index        int    `json:"index"`
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason,omitempty"`
}

// Delta 流式增量消息
type Delta struct {
	Role      string     `json:"role,omitempty"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Usage token使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Tool 工具定义
type Tool struct {
	Type     string       `json:"type"`     // 工具类型，当前只支持"function"
	Function ToolFunction `json:"function"` // 函数配置
}

// ToolFunction 工具函数定义
type ToolFunction struct {
	Name        string      `json:"name"`                  // 函数名（必需）
	Description string      `json:"description,omitempty"` // 函数描述
	Parameters  interface{} `json:"parameters,omitempty"`  // 函数参数，JSON Schema格式
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string           `json:"id"`       // 工具调用ID
	Type     string           `json:"type"`     // 工具类型
	Function ToolCallFunction `json:"function"` // 函数信息
}

// ToolCallFunction 工具调用函数信息
type ToolCallFunction struct {
	Name      string `json:"name"`      // 函数名
	Arguments string `json:"arguments"` // 函数参数，JSON字符串
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	OwnedBy string `json:"owned_by"`
}

// ModelsResponse 模型列表响应
type ModelsResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

// ErrorResponse API 错误响应
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}
