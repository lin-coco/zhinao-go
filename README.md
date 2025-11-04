# 360智脑 Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lin-coco/zhinao-go)](https://goreportcard.com/report/github.com/lin-coco/zhinao-go)
[![GoDoc](https://godoc.org/github.com/lin-coco/zhinao-go?status.svg)](https://godoc.org/github.com/lin-coco/zhinao-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

360智脑 Go 语言 SDK，提供简洁、类型安全的 API 接口来访问360智脑的 AI 能力。

## ✨ 特性

- 🚀 **简单易用** - 直观的 API 和链式 Builder 模式
- 🔄 **流式支持** - 完整的流式聊天补全
- ⚡ **高性能** - 连接复用和智能重试
- 🛡️ **类型安全** - 完整的类型定义
- 🔧 **灵活配置** - 函数式选项模式
- 📝 **完善文档** - 详细的文档和示例

## 📦 安装

```bash
go get github.com/lin-coco/zhinao-go
```

要求 Go 版本 >= 1.18

## 🚀 快速开始

### 基础使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/lin-coco/zhinao-go"
)

func main() {
    // 创建客户端（自动从环境变量 ZHINAO_API_KEY 读取）
    client, err := zhinao.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 使用 Builder 构建请求（推荐）
    req := zhinao.NewChatBuilder().
        SetModel(zhinao.Model360GPTTurbo).
        AddUserMessage("你好，请介绍一下360智脑").
        Build()

    // 发送请求
    resp, err := client.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 环境变量配置

```bash
# 设置 API Key
export ZHINAO_API_KEY="your-api-key"
```

### 流式响应

```go
stream, err := client.Chat.CreateCompletionStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(resp.Choices[0].Delta.Content)
}
```

### 自定义配置

```go
client, err := zhinao.NewClient(
    "your-api-key",
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithRetry(5, 2*time.Second),
)
```

## 📖 核心 API

### 客户端创建

```go
// 方式1: 从环境变量创建（推荐）
client, err := zhinao.NewClientFromEnv()

// 方式2: 直接传入 API Key
client, err := zhinao.NewClient("your-api-key")

// 方式3: 带配置选项
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithRetry(5, 2*time.Second),
)
```

### Builder 模式（推荐）

```go
req := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddSystemMessage("你是一个助手").
    AddUserMessage("你好").
    SetTemperature(0.7).
    SetMaxTokens(1000).
    Build()
```

### 聊天补全

```go
// 非流式
resp, err := client.Chat.CreateCompletion(ctx, req)

// 流式
stream, err := client.Chat.CreateCompletionStream(ctx, req)
```

### 模型列表

```go
models, err := client.Models.List(ctx)
```

## 📚 示例代码

查看 [examples](./examples) 目录获取完整示例：

| 示例 | 说明 | 运行 |
|-----|------|------|
| [chat-completion](./examples/chat-completion/) | 基础聊天补全 | `go run examples/chat-completion/main.go` |
| [chatbot](./examples/chatbot/) | 交互式聊天机器人 | `go run examples/chatbot/main.go` |
| [stream-chat](./examples/stream-chat/) | 流式响应 | `go run examples/stream-chat/main.go` |
| [chat-with-tools](./examples/chat-with-tools/) | 工具调用 | `go run examples/chat-with-tools/main.go` |
| [chat-with-builder](./examples/chat-with-builder/) | Builder 模式 | `go run examples/chat-with-builder/main.go` |

详细说明请查看 [examples/README.md](./examples/README.md)

## 📘 文档

- **[完整指南](./docs/GUIDE.md)** - 详细的使用指南，包含架构设计、高级用法、测试策略等
- **[SDK 对比](./docs/COMPARISON.md)** - 与其他流行 SDK 的对比分析
- **[文档导航](./docs/README.md)** - 文档结构说明

### 快速链接

- [架构设计](./docs/GUIDE.md#架构设计)
- [错误处理](./docs/GUIDE.md#4-错误处理)
- [测试指南](./docs/GUIDE.md#测试指南)
- [最佳实践](./docs/GUIDE.md#最佳实践)
- [扩展开发](./docs/GUIDE.md#扩展开发)

## 🛠️ 配置选项

| 选项 | 说明 | 示例 |
|-----|------|------|
| `WithTimeout` | 设置请求超时 | `WithTimeout(30*time.Second)` |
| `WithRetry` | 设置重试策略 | `WithRetry(5, 2*time.Second)` |
| `WithBaseURL` | 自定义 API 地址 | `WithBaseURL("https://api.360.cn/v1")` |
| `WithUserAgent` | 自定义 User-Agent | `WithUserAgent("MyApp/1.0")` |
| `WithHeaders` | 自定义请求头 | `WithHeaders(map[string]string{...})` |

## 🤝 贡献

欢迎贡献！请查看 [docs/GUIDE.md - 贡献指南](./docs/GUIDE.md#贡献指南)

提交 PR 前请确保：
- ✅ 测试通过 (`make test`)
- ✅ 代码检查通过 (`make lint`)
- ✅ 添加了测试用例
- ✅ 更新了文档

## 📄 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件

## 🔗 相关链接

- [360智脑官网](https://ai.360.com)
- [API 文档](https://ai.360.com/platform/docs/overview)
- [问题反馈](https://github.com/lin-coco/zhinao-go/issues)

## 🙏 致谢

本项目参考了以下优秀开源项目的设计：
- [go-openai](https://github.com/sashabaranov/go-openai)
- [go-moonshot](https://github.com/northes/go-moonshot)
- [deepseek-go](https://github.com/cohesion-org/deepseek-go)

详细对比分析请查看 [docs/COMPARISON.md](./docs/COMPARISON.md)
