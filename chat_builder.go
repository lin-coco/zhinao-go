package zhinao

// ChatCompletionBuilder 聊天请求构建器
//
// 提供链式调用的方式构建聊天请求，简化复杂请求的创建过程
type ChatCompletionBuilder struct {
	req *ChatCompletionRequest
}

// NewChatCompletionBuilder 创建新的聊天请求构建器
//
// 示例:
//
//	builder := zhinao.NewChatCompletionBuilder().
//	    SetModel("360gpt-turbo").
//	    AddUserMessage("你好").
//	    SetTemperature(0.7).
//	    Build()
func NewChatCompletionBuilder() *ChatCompletionBuilder {
	return &ChatCompletionBuilder{
		req: &ChatCompletionRequest{
			Messages: make([]ChatCompletionMessage, 0),
		},
	}
}

// SetModel 设置模型
func (b *ChatCompletionBuilder) SetModel(model string) *ChatCompletionBuilder {
	b.req.Model = model
	return b
}

// AddMessage 添加消息
func (b *ChatCompletionBuilder) AddMessage(role, content string) *ChatCompletionBuilder {
	b.req.Messages = append(b.req.Messages, ChatCompletionMessage{
		Role:    role,
		Content: content,
	})
	return b
}

// AddSystemMessage 添加系统消息
//
// 系统消息用于设置对话的上下文和行为
func (b *ChatCompletionBuilder) AddSystemMessage(content string) *ChatCompletionBuilder {
	return b.AddMessage("system", content)
}

// AddUserMessage 添加用户消息
func (b *ChatCompletionBuilder) AddUserMessage(content string) *ChatCompletionBuilder {
	return b.AddMessage("user", content)
}

// AddAssistantMessage 添加助手消息
//
// 助手消息通常用于多轮对话中提供历史回复
func (b *ChatCompletionBuilder) AddAssistantMessage(content string) *ChatCompletionBuilder {
	return b.AddMessage("assistant", content)
}

// AddMessages 批量添加消息
func (b *ChatCompletionBuilder) AddMessages(messages []ChatCompletionMessage) *ChatCompletionBuilder {
	b.req.Messages = append(b.req.Messages, messages...)
	return b
}

// SetTemperature 设置温度参数
//
// 温度控制输出的随机性，范围 [0, 1]
// - 较低的值（如 0.2）使输出更确定和聚焦
// - 较高的值（如 0.8）使输出更随机和创造性
func (b *ChatCompletionBuilder) SetTemperature(temperature float64) *ChatCompletionBuilder {
	b.req.Temperature = temperature
	return b
}

// SetMaxTokens 设置最大生成 token 数
//
// 限制生成文本的长度，范围 [1, 2048]
func (b *ChatCompletionBuilder) SetMaxTokens(maxTokens int) *ChatCompletionBuilder {
	b.req.MaxTokens = maxTokens
	return b
}

// SetTopP 设置核采样参数
//
// 控制生成文本的多样性，范围 [0, 1]
// 建议只修改 temperature 或 top_p 其中之一
func (b *ChatCompletionBuilder) SetTopP(topP float64) *ChatCompletionBuilder {
	b.req.TopP = topP
	return b
}

// SetTopK 设置 TopK 参数
//
// 限制每次采样时考虑的 token 数量，范围 [0, 1024]
func (b *ChatCompletionBuilder) SetTopK(topK int) *ChatCompletionBuilder {
	b.req.TopK = topK
	return b
}

// SetRepetitionPenalty 设置重复惩罚
//
// 减少重复内容的生成，范围 [1, 2]
func (b *ChatCompletionBuilder) SetRepetitionPenalty(penalty float64) *ChatCompletionBuilder {
	b.req.RepetitionPenalty = penalty
	return b
}

// SetNumBeams 设置 beam search 数量
//
// 使用 beam search 提高生成质量，范围 [1, 5]
func (b *ChatCompletionBuilder) SetNumBeams(numBeams int) *ChatCompletionBuilder {
	b.req.NumBeams = numBeams
	return b
}

// AddTool 添加工具
//
// 工具允许模型调用外部函数
func (b *ChatCompletionBuilder) AddTool(tool Tool) *ChatCompletionBuilder {
	b.req.Tools = append(b.req.Tools, tool)
	return b
}

// AddTools 批量添加工具
func (b *ChatCompletionBuilder) AddTools(tools []Tool) *ChatCompletionBuilder {
	b.req.Tools = append(b.req.Tools, tools...)
	return b
}

// SetToolChoice 设置工具选择策略
//
// 可以是字符串 "none"、"auto" 或具体的工具选择对象
func (b *ChatCompletionBuilder) SetToolChoice(toolChoice interface{}) *ChatCompletionBuilder {
	b.req.ToolChoice = toolChoice
	return b
}

// SetUser 设置用户标识
//
// 用于追踪和分析用户行为
func (b *ChatCompletionBuilder) SetUser(user string) *ChatCompletionBuilder {
	b.req.User = user
	return b
}

// Build 构建聊天请求
//
// 返回构建好的请求对象
func (b *ChatCompletionBuilder) Build() *ChatCompletionRequest {
	return b.req
}

// BuildAndValidate 构建并验证聊天请求
//
// 返回构建好的请求对象，如果验证失败则返回错误
func (b *ChatCompletionBuilder) BuildAndValidate() (*ChatCompletionRequest, error) {
	if err := b.req.Validate(); err != nil {
		return nil, err
	}
	return b.req, nil
}
