package zhinao

// ChatBuilder 聊天请求构建器
//
// 提供链式调用的方式构建聊天请求，简化复杂请求的创建过程
type ChatBuilder struct {
	req *ChatRequest
}

// NewChatBuilder 创建新的聊天请求构建器
//
// 示例:
//
//	builder := zhinao.NewChatBuilder().
//	    SetModel("360gpt-turbo").
//	    AddUserMessage("你好").
//	    SetTemperature(0.7).
//	    Build()
func NewChatBuilder() *ChatBuilder {
	return &ChatBuilder{
		req: &ChatRequest{
			Messages: make([]Message, 0),
		},
	}
}

// SetModel 设置模型
func (b *ChatBuilder) SetModel(model string) *ChatBuilder {
	b.req.Model = model
	return b
}

// AddMessage 添加消息
func (b *ChatBuilder) AddMessage(role, content string) *ChatBuilder {
	b.req.Messages = append(b.req.Messages, Message{
		Role:    role,
		Content: content,
	})
	return b
}

// AddSystemMessage 添加系统消息
//
// 系统消息用于设置对话的上下文和行为
func (b *ChatBuilder) AddSystemMessage(content string) *ChatBuilder {
	return b.AddMessage("system", content)
}

// AddUserMessage 添加用户消息
func (b *ChatBuilder) AddUserMessage(content string) *ChatBuilder {
	return b.AddMessage("user", content)
}

// AddAssistantMessage 添加助手消息
//
// 助手消息通常用于多轮对话中提供历史回复
func (b *ChatBuilder) AddAssistantMessage(content string) *ChatBuilder {
	return b.AddMessage("assistant", content)
}

// AddMessages 批量添加消息
func (b *ChatBuilder) AddMessages(messages []Message) *ChatBuilder {
	b.req.Messages = append(b.req.Messages, messages...)
	return b
}

// SetTemperature 设置温度参数
//
// 温度控制输出的随机性，范围 [0, 1]
// - 较低的值（如 0.2）使输出更确定和聚焦
// - 较高的值（如 0.8）使输出更随机和创造性
func (b *ChatBuilder) SetTemperature(temperature float64) *ChatBuilder {
	b.req.Temperature = temperature
	return b
}

// SetMaxTokens 设置最大生成 token 数
//
// 限制生成文本的长度，范围 [1, 2048]
func (b *ChatBuilder) SetMaxTokens(maxTokens int) *ChatBuilder {
	b.req.MaxTokens = maxTokens
	return b
}

// SetTopP 设置核采样参数
//
// 控制生成文本的多样性，范围 [0, 1]
// 建议只修改 temperature 或 top_p 其中之一
func (b *ChatBuilder) SetTopP(topP float64) *ChatBuilder {
	b.req.TopP = topP
	return b
}

// SetTopK 设置 TopK 参数
//
// 限制每次采样时考虑的 token 数量，范围 [0, 1024]
func (b *ChatBuilder) SetTopK(topK int) *ChatBuilder {
	b.req.TopK = topK
	return b
}

// SetRepetitionPenalty 设置重复惩罚
//
// 减少重复内容的生成，范围 [1, 2]
func (b *ChatBuilder) SetRepetitionPenalty(penalty float64) *ChatBuilder {
	b.req.RepetitionPenalty = penalty
	return b
}

// SetNumBeams 设置 beam search 数量
//
// 使用 beam search 提高生成质量，范围 [1, 5]
func (b *ChatBuilder) SetNumBeams(numBeams int) *ChatBuilder {
	b.req.NumBeams = numBeams
	return b
}

// AddTool 添加工具
//
// 工具允许模型调用外部函数
func (b *ChatBuilder) AddTool(tool Tool) *ChatBuilder {
	b.req.Tools = append(b.req.Tools, tool)
	return b
}

// AddTools 批量添加工具
func (b *ChatBuilder) AddTools(tools []Tool) *ChatBuilder {
	b.req.Tools = append(b.req.Tools, tools...)
	return b
}

// SetToolChoice 设置工具选择策略
//
// 可以是字符串 "none"、"auto" 或具体的工具选择对象
func (b *ChatBuilder) SetToolChoice(toolChoice interface{}) *ChatBuilder {
	b.req.ToolChoice = toolChoice
	return b
}

// SetUser 设置用户标识
//
// 用于追踪和分析用户行为
func (b *ChatBuilder) SetUser(user string) *ChatBuilder {
	b.req.User = user
	return b
}

// Build 构建聊天请求
//
// 返回构建好的请求对象
func (b *ChatBuilder) Build() *ChatRequest {
	return b.req
}

// BuildAndValidate 构建并验证聊天请求
//
// 返回构建好的请求对象，如果验证失败则返回错误
func (b *ChatBuilder) BuildAndValidate() (*ChatRequest, error) {
	if err := b.req.Validate(); err != nil {
		return nil, err
	}
	return b.req, nil
}
