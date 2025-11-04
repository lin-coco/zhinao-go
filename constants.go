package zhinao

// 360智脑自有模型
const (
	// Model360GPTTurbo 360GPT Turbo 模型
	Model360GPTTurbo = "360gpt-turbo"

	// Model360GPTPro 360GPT Pro 模型
	Model360GPTPro = "360gpt-pro"

	// Model360GPT2Pro 360GPT2 Pro 模型
	Model360GPT2Pro = "360gpt2-pro"

	// Model360Zhinao3Pro 360智脑3 Pro 模型
	Model360Zhinao3Pro = "360zhinao3-pro"
)

// DeepSeek 模型
const (
	// ModelDeepSeekV3 DeepSeek V3 模型
	ModelDeepSeekV3 = "deepseek/deepseek-chat"

	// ModelDeepSeekR1 DeepSeek R1 推理模型
	ModelDeepSeekR1 = "deepseek/deepseek-reasoner"
)

// 阿里通义千问模型
const (
	// ModelQwenMax 通义千问 Max 模型
	ModelQwenMax = "alibaba/qwen-max"

	// ModelQwenPlus 通义千问 Plus 模型
	ModelQwenPlus = "alibaba/qwen-plus"

	// ModelQwenTurbo 通义千问 Turbo 模型
	ModelQwenTurbo = "alibaba/qwen-turbo"

	// ModelQwenVLMax 通义千问视觉 Max 模型
	ModelQwenVLMax = "alibaba/qwen-vl-max"
)

// 百度文心模型
const (
	// ModelErnie4Turbo 文心一言 4.0 Turbo 模型
	ModelErnie4Turbo = "baidu/ernie-4.5-turbo-128k"

	// ModelErnie35 文心一言 3.5 模型
	ModelErnie35 = "baidu/ernie-3.5-8k"
)

// 字节豆包模型
const (
	// ModelDoubaoProMax 豆包 Pro Max 模型
	ModelDoubaoProMax = "volcengine/doubao-seed-1-6"

	// ModelDoubaoPro 豆包 Pro 模型
	ModelDoubaoPro = "Doubao-pro-32k"
)

// MiniMax 模型
const (
	// ModelMiniMaxText01 MiniMax Text-01 模型
	ModelMiniMaxText01 = "minimax/MiniMax-Text-01"

	// ModelMiniMaxAbab65 MiniMax abab6.5 模型
	ModelMiniMaxAbab65 = "minimax/abab6.5-chat"
)

// 智谱 GLM 模型
const (
	// ModelGLM4Plus 智谱 GLM-4 Plus 模型
	ModelGLM4Plus = "bigmodel/glm-4-plus"

	// ModelGLM4 智谱 GLM-4 模型
	ModelGLM4 = "bigmodel/glm-4"

	// ModelGLM4V 智谱 GLM-4V 视觉模型
	ModelGLM4V = "bigmodel/glm-4v"
)

// 月之暗面 Moonshot 模型
const (
	// ModelMoonshotV1 Moonshot V1 128K 模型
	ModelMoonshotV1 = "moonshot/moonshot-v1-128k"

	// ModelKimiK2 Kimi K2 模型
	ModelKimiK2 = "moonshot/kimi-k2-0905-preview"
)

// OpenAI 模型
const (
	// ModelGPT4o GPT-4o 模型
	ModelGPT4o = "openai/gpt-4o"

	// ModelGPT4oMini GPT-4o Mini 模型
	ModelGPT4oMini = "gpt-4o-mini"

	// ModelO1 OpenAI o1 推理模型
	ModelO1 = "openai/o1"
)

// Anthropic Claude 模型
const (
	// ModelClaudeSonnet4 Claude Sonnet 4 模型
	ModelClaudeSonnet4 = "anthropic/claude-sonnet-4"

	// ModelClaude35Sonnet Claude 3.5 Sonnet 模型
	ModelClaude35Sonnet = "anthropic/claude-3.5-sonnet"
)

// Google Gemini 模型
const (
	// ModelGemini2Flash Gemini 2.0 Flash 模型
	ModelGemini2Flash = "google/gemini-2.0-flash"

	// ModelGemini25Pro Gemini 2.5 Pro 模型
	ModelGemini25Pro = "google/gemini-2.5-pro"
)

// 腾讯混元模型
const (
	// ModelHunyuanTurbo 混元 Turbo 模型
	ModelHunyuanTurbo = "tencent/hunyuan-turbo"

	// ModelHunyuanPro 混元 Pro 模型
	ModelHunyuanPro = "tencent/hunyuan-pro"
)

// 商汤日日新模型
const (
	// ModelSenseChat5 日日新 5.0 模型
	ModelSenseChat5 = "sensetime/SenseChat-5"

	// ModelSenseChatTurbo 日日新 Turbo 模型
	ModelSenseChatTurbo = "sensetime/SenseChat-Turbo"
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
