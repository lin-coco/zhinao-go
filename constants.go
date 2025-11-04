package zhinao

// 模型常量
const (
	// Model360GPTTurbo 360GPT Turbo 模型
	Model360GPTTurbo = "360gpt-turbo"

	// Model360GPTPro 360GPT Pro 模型
	Model360GPTPro = "360gpt-pro"

	// Model360GPT2Pro 360GPT2 Pro 模型
	Model360GPT2Pro = "360gpt2-pro"
)

// 角色常量
const (
	// RoleSystem 系统角色
	RoleSystem = "system"

	// RoleUser 用户角色
	RoleUser = "user"

	// RoleAssistant 助手角色
	RoleAssistant = "assistant"

	// RoleTool 工具角色
	RoleTool = "tool"
)

// 工具类型常量
const (
	// ToolTypeFunction 函数工具类型
	ToolTypeFunction = "function"
)

// 完成原因常量
const (
	// FinishReasonStop 正常停止
	FinishReasonStop = "stop"

	// FinishReasonLength 达到最大长度
	FinishReasonLength = "length"

	// FinishReasonToolCalls 工具调用
	FinishReasonToolCalls = "tool_calls"

	// FinishReasonContentFilter 内容过滤
	FinishReasonContentFilter = "content_filter"
)
