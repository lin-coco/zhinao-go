package zhinao

// ============================================================================
// Chat Completion Types
// ============================================================================

// ChatCompletionMessage 聊天消息
type ChatCompletionMessage struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // 消息内容
	Name       string     `json:"name,omitempty"`         // 消息发送者名称（可选）
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 工具调用列表（仅助手消息）
	ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID（仅工具消息）
}

// ChatCompletionRequest 聊天请求
type ChatCompletionRequest struct {
	Model             string                  `json:"model"`                        // 模型名称（必需）
	Messages          []ChatCompletionMessage `json:"messages"`                     // 消息列表（必需）
	Stream            bool                    `json:"stream,omitempty"`             // 是否流式输出
	Temperature       float64                 `json:"temperature,omitempty"`        // 温度参数 [0, 1]，默认0.9
	MaxTokens         int                     `json:"max_tokens,omitempty"`         // 最大生成token数 [1, 2048]，默认2048
	TopP              float64                 `json:"top_p,omitempty"`              // 核采样参数 [0, 1]，默认0.5
	TopK              int                     `json:"top_k,omitempty"`              // TopK参数 [0, 1024]，默认0
	RepetitionPenalty float64                 `json:"repetition_penalty,omitempty"` // 重复惩罚 [1, 2]，默认1.05
	NumBeams          int                     `json:"num_beams,omitempty"`          // beam search数量 [1, 5]，默认1
	Tools             []Tool                  `json:"tools,omitempty"`              // 工具列表
	ToolChoice        interface{}             `json:"tool_choice,omitempty"`        // 工具选择，可以是字符串或结构体
	User              string                  `json:"user,omitempty"`               // 用户标识
}

// Validate 验证请求参数
func (r *ChatCompletionRequest) Validate() error {
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

// ChatCompletionResponse 聊天响应
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

// ChatCompletionChoice 响应选择项
type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionStreamResponse 流式聊天响应
type ChatCompletionStreamResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []ChatCompletionStreamChoice `json:"choices"`
}

// ChatCompletionStreamChoice 流式响应选择项
type ChatCompletionStreamChoice struct {
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

// ============================================================================
// Embedding Types
// ============================================================================

// EmbeddingRequest 向量生成请求
type EmbeddingRequest struct {
	// Model 模型类型（必填）
	Model string `json:"model"`

	// Input 批量生成，每一项是需要生成向量的内容（必填）
	Input []string `json:"input"`

	// User 标记业务方用户 id，便于业务方区分不同用户（可选）
	User string `json:"user,omitempty"`
}

// Validate 验证请求参数
func (req *EmbeddingRequest) Validate() error {
	if req.Model == "" {
		return ErrInvalidModel
	}

	if len(req.Input) == 0 {
		return ErrEmptyInput
	}

	return nil
}

// EmbeddingResponse 向量生成响应
type EmbeddingResponse struct {
	// Data 返回结果，一个结构体数组，每个元素对应 input 的每一项输入
	Data []Embedding `json:"data"`

	// Model 本次调用使用的模型名
	Model string `json:"model"`

	// Object 对象类型
	Object string `json:"object"`

	// Usage token 消耗量
	Usage EmbeddingUsage `json:"usage"`
}

// Embedding 向量数据
type Embedding struct {
	// Embedding 返回的向量结果，是一个浮点数数组
	Embedding []float64 `json:"embedding"`

	// Object 对象类型
	Object string `json:"object"`

	// Index 代表本次结果在 data 里的下标值
	Index int `json:"index"`
}

// EmbeddingUsage token 使用量
type EmbeddingUsage struct {
	// PromptTokens 输入 token 消耗量
	PromptTokens int `json:"prompt_tokens"`

	// TotalTokens 总 token 消耗量
	TotalTokens int `json:"total_tokens"`
}

// ============================================================================
// Image Types
// ============================================================================

// ImageStyle 图像风格类型
type ImageStyle string

const (
	// ImageStylePapercut 剪纸风格
	ImageStylePapercut ImageStyle = "papercut"

	// ImageStyleRealistic 写实风格
	ImageStyleRealistic ImageStyle = "realistic"

	// ImageStyleCartoon 卡通风格
	ImageStyleCartoon ImageStyle = "cartoon"

	// ImageStyleCG CG风格
	ImageStyleCG ImageStyle = "CG"
)

// ImageRequest 文本生成图像请求
type ImageRequest struct {
	// Model 模型类型（必填）
	Model string `json:"model"`

	// Style 生成的图像风格（必填）
	// 可选风格：papercut(剪纸)、realistic(写实)、cartoon(卡通)、CG
	Style ImageStyle `json:"style"`

	// Prompt 用于生成图像的提示信息（必填）
	// 长度应小于等于 500 个字符
	Prompt string `json:"prompt"`

	// NegativePrompt 用于生成图像的负向提示信息（可选）
	// 即不希望出现的元素，长度应小于等于 500 个字符
	NegativePrompt string `json:"negative_prompt,omitempty"`

	// Width 生成的图像宽度（可选）
	// 取值应大于等于 8 小于等于 2048，默认是 512
	Width int `json:"width,omitempty"`

	// Height 生成的图像高度（可选）
	// 取值应大于等于 8 小于等于 2048，默认是 512
	Height int `json:"height,omitempty"`

	// Samples 需要生成的图片数量（可选）
	// 取值应大于等于 1 小于等于 4，默认是 1
	Samples int `json:"samples,omitempty"`

	// NumInferenceSteps 采样步数（可选）
	// 取值应大于等于 20 小于等于 50，默认是 25
	NumInferenceSteps int `json:"num_inference_steps,omitempty"`

	// Seed 种子（可选）
	// 正整数，默认随机，用于控制生成图片的随机性
	Seed int `json:"seed,omitempty"`

	// GuidanceScale 提示词强度（可选）
	// 取值应大于等于 7.5 小于等于 15，默认是 7.5
	GuidanceScale float64 `json:"guidance_scale,omitempty"`

	// EnhancePrompt 是否进行 prompt 润色（可选）
	// 默认是 false
	EnhancePrompt bool `json:"enhance_prompt,omitempty"`
}

// Validate 验证请求参数
func (r *ImageRequest) Validate() error {
	if r.Model == "" {
		return ErrInvalidModel
	}

	if r.Style == "" {
		return &ValidationError{Message: "style cannot be empty"}
	}

	if r.Prompt == "" {
		return &ValidationError{Message: "prompt cannot be empty"}
	}

	if len(r.Prompt) > 500 {
		return &ValidationError{Message: "prompt length must be <= 500"}
	}

	if r.NegativePrompt != "" && len(r.NegativePrompt) > 500 {
		return &ValidationError{Message: "negative_prompt length must be <= 500"}
	}

	if r.Width != 0 && (r.Width < 8 || r.Width > 2048) {
		return &ValidationError{Message: "width must be between 8 and 2048"}
	}

	if r.Height != 0 && (r.Height < 8 || r.Height > 2048) {
		return &ValidationError{Message: "height must be between 8 and 2048"}
	}

	if r.Samples != 0 && (r.Samples < 1 || r.Samples > 4) {
		return &ValidationError{Message: "samples must be between 1 and 4"}
	}

	if r.NumInferenceSteps != 0 && (r.NumInferenceSteps < 20 || r.NumInferenceSteps > 50) {
		return &ValidationError{Message: "num_inference_steps must be between 20 and 50"}
	}

	if r.GuidanceScale != 0 && (r.GuidanceScale < 7.5 || r.GuidanceScale > 15) {
		return &ValidationError{Message: "guidance_scale must be between 7.5 and 15"}
	}

	return nil
}

// ImageResponse 文本生成图像响应
type ImageResponse struct {
	// Status 状态
	Status string `json:"status"`

	// GenerationTime 耗时，以秒为单位
	GenerationTime int `json:"generationTime"`

	// Output 生成的图片链接数组
	// 如果生成数少于指定的数，一般是因为部分图片涉及敏感被过滤
	Output []string `json:"output"`

	// Meta 元数据
	Meta ImageMeta `json:"meta"`
}

// ImageMeta 图像元数据
type ImageMeta struct {
	// H 图像高度
	H int `json:"H"`

	// W 图像宽度
	W int `json:"W"`

	// GuidanceScale 提示词强度
	GuidanceScale float64 `json:"guidance_scale"`

	// NSamples 一次请求生成的图片数
	NSamples int `json:"n_samples"`

	// NegativePrompt 用于图像生成的负向 prompt
	NegativePrompt string `json:"negative_prompt"`

	// Prompt 用于图像生成的 prompt
	Prompt string `json:"prompt"`

	// Seed 种子
	Seed int `json:"seed"`

	// Steps 采样步数
	Steps int `json:"steps"`
}

// ============================================================================
// Model Types
// ============================================================================

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

// ============================================================================
// Error Types
// ============================================================================

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
