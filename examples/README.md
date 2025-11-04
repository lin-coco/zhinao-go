# 360智脑 Go SDK 示例

本目录包含了 360智脑 Go SDK 的各种使用示例。所有示例都需要设置 `ZHINAO_API_KEY` 环境变量。

## 设置 API Key

```bash
export ZHINAO_API_KEY="your-api-key-here"
```

## 示例列表

### 1. chat-completion - 基础聊天补全

最简单的聊天补全示例，展示如何发送单个消息并获取 AI 回复。

```bash
cd examples/chat-completion
go run main.go
```

**适用场景**: 
- 快速开始使用 SDK
- 单次问答
- API 功能测试

---

### 2. chatbot - 交互式聊天机器人

实现一个可以持续对话的聊天机器人，支持多轮对话和上下文记忆。

```bash
cd examples/chatbot
go run main.go
```

**功能特点**:
- 支持多轮对话
- 维护对话历史
- 交互式命令行界面
- 输入 `exit` 退出

**适用场景**:
- 构建聊天应用
- 客服机器人
- 智能助手

---

### 3. stream-chat - 流式响应

演示如何使用流式 API 实时接收 AI 响应，适合需要即时反馈的场景。

```bash
cd examples/stream-chat
go run main.go
```

**功能特点**:
- 实时流式输出
- 降低首字延迟
- 更好的用户体验

**适用场景**:
- 长文本生成
- 需要实时反馈的应用
- 提升用户体验

---

### 4. chat-with-tools - 工具调用 (Function Calling)

展示如何使用工具调用功能，让 AI 能够调用外部函数来获取信息或执行操作。

```bash
cd examples/chat-with-tools
go run main.go
```

**功能特点**:
- 定义工具函数
- AI 自主决定是否调用工具
- 基于工具结果生成回答

**适用场景**:
- 需要访问外部数据（天气、股票、数据库等）
- 执行特定操作（发送邮件、创建任务等）
- 增强 AI 能力

---

### 5. chat-with-builder - Builder 模式

演示如何使用链式 Builder 模式构建请求，提供更流畅的 API 使用体验。

```bash
cd examples/chat-with-builder
go run main.go
```

**功能特点**:
- 链式调用
- 类型安全
- 代码可读性强
- 灵活配置参数

**适用场景**:
- 需要配置多个参数
- 提高代码可维护性
- 推荐的最佳实践

---

### 6. list-models - 模型列表

展示如何获取可用模型列表和查询特定模型信息。

```bash
cd examples/list-models
go run main.go
```

**功能特点**:
- 获取所有可用模型
- 查询特定模型详情
- 显示模型基本信息

**适用场景**:
- 了解可用模型
- 动态选择模型
- 模型信息查询

---

### 7. text2img - 图像生成

展示如何使用文本生成图像功能，支持多种风格和参数配置。

```bash
cd examples/text2img
go run main.go
```

**功能特点**:
- 文本生成图像
- 多种风格选择（写实、卡通、剪纸、CG）
- 支持负向提示词
- 批量生成
- 自定义参数（尺寸、步数、种子等）

**适用场景**:
- AI 绘画应用
- 内容创作工具
- 图片素材生成
- 创意设计辅助

---

## 运行所有示例

你可以使用以下命令快速运行某个示例：

```bash
# 运行聊天补全示例
go run examples/chat-completion/main.go

# 运行聊天机器人
go run examples/chatbot/main.go

# 运行流式聊天
go run examples/stream-chat/main.go

# 运行工具调用示例
go run examples/chat-with-tools/main.go

# 运行 Builder 示例
go run examples/chat-with-builder/main.go

# 运行模型列表示例
go run examples/list-models/main.go

# 运行图像生成示例
go run examples/text2img/main.go
```

## 常见问题

### 1. 如何获取 API Key?

请访问 [360智脑官网](https://ai.360.com) 注册账号并获取 API Key。

### 2. 示例运行失败怎么办?

- 确认已设置 `ZHINAO_API_KEY` 环境变量
- 检查 API Key 是否有效
- 确认网络连接正常
- 查看错误信息进行排查

### 3. 如何自定义示例?

所有示例都是独立的 Go 程序，你可以：
- 修改消息内容
- 调整模型参数
- 添加更多功能
- 集成到你的应用中

## 更多资源

- [完整文档](../README.md)
- [API 参考](../docs/ARCHITECTURE.md)
- [最佳实践](../docs/COMPARISON.md)
- [测试指南](../docs/TESTING.md)

## 贡献

欢迎提交更多示例！如果你有好的使用案例，请：
1. Fork 本项目
2. 创建新的示例目录
3. 添加清晰的注释和说明
4. 提交 Pull Request
