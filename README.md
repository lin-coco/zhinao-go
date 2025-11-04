# 360智脑 Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lin-coco/zhinao-go)](https://goreportcard.com/report/github.com/lin-coco/zhinao-go)
[![GoDoc](https://godoc.org/github.com/lin-coco/zhinao-go?status.svg)](https://godoc.org/github.com/lin-coco/zhinao-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

360智脑 Go 语言 SDK，提供简洁、类型安全且功能强大的 API 接口来访问 360 智脑的 AI 能力。

## ✨ 核心特性

- **简洁易用** - 直观的 API 设计，支持环境变量便捷初始化
- **Builder 模式** - 流畅的链式调用，轻松构建复杂请求
- **流式响应** - 实时接收 AI 生成内容
- **自动重试** - 内置智能重试机制，指数退避策略
- **类型安全** - 完整的类型定义和编译时检查
- **生产就绪** - 完善的错误处理和测试覆盖

## 📦 安装

```bash
go get github.com/lin-coco/zhinao-go
```

**要求**: Go >= 1.18

## 🚀 快速开始

### 获取 API Key

访问 [360智脑开放平台](https://ai.360.com) 获取 API Key

### 设置环境变量

```bash
export ZHINAO_API_KEY="your-api-key-here"
```

### 基础示例

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

    // 使用 Builder 构建请求
    req := zhinao.NewChatBuilder().
        SetModel(zhinao.Model360GPTTurbo).
        AddUserMessage("用一句话介绍Go语言的特点").
        Build()

    // 发送请求
    resp, err := client.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 流式响应示例

```go
// 构建请求，启用流式
req := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddUserMessage("写一首关于秋天的诗").
    SetStream(true).
    Build()

// 创建流式响应
stream, err := client.Chat.CreateCompletionStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

// 实时接收并打印内容
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

## 📚 完整示例

查看 [examples](./examples) 目录获取完整可运行的示例：

| 示例 | 说明 | 运行命令 |
|-----|------|---------|
| [chat-completion](./examples/chat-completion/) | 基础聊天补全 | `go run examples/chat-completion/main.go` |
| [chatbot](./examples/chatbot/) | 交互式聊天机器人 | `go run examples/chatbot/main.go` |
| [stream-chat](./examples/stream-chat/) | 流式响应 | `go run examples/stream-chat/main.go` |
| [chat-with-tools](./examples/chat-with-tools/) | 工具调用 | `go run examples/chat-with-tools/main.go` |
| [chat-with-builder](./examples/chat-with-builder/) | Builder 模式 | `go run examples/chat-with-builder/main.go` |
| [list-models](./examples/list-models/) | 模型列表 | `go run examples/list-models/main.go` |
| [text2img](./examples/text2img/) | 文本生成图像 | `go run examples/text2img/main.go` |
| [embeddings](./examples/embeddings/) | 向量生成 | `go run examples/embeddings/main.go` |

详细说明请查看 [examples/README.md](./examples/README.md)

## 🔧 配置选项

```go
client, err := zhinao.NewClient(
    "your-api-key",
    zhinao.WithTimeout(30*time.Second),           // 设置超时
    zhinao.WithRetry(5, 2*time.Second),           // 重试配置
    zhinao.WithBaseURL("https://api.360.cn/v1"),  // 自定义 API 地址
    zhinao.WithUserAgent("MyApp/1.0"),            // 自定义 User-Agent
)
```

## 📘 文档

- **[完整指南](./docs/GUIDE.md)** - 详细的使用指南，包含架构设计、高级用法、测试策略等
- **[SDK 对比](./docs/COMPARISON.md)** - 与其他流行 Go SDK 的设计对比分析

### 快速导航

- [架构设计](./docs/GUIDE.md#架构设计)
- [核心功能](./docs/GUIDE.md#核心功能)
- [测试指南](./docs/GUIDE.md#测试指南)
- [最佳实践](./docs/GUIDE.md#最佳实践)

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](./docs/GUIDE.md#贡献指南)

### 开发命令

```bash
# 运行测试
make test

# 代码检查
make lint

# 生成覆盖率报告
make test-coverage

# 格式化代码
make fmt
```

## 🌟 设计理念

360智脑 Go SDK 的设计参考了多个优秀的开源项目：

- **[go-openai](https://github.com/sashabaranov/go-openai)** - 流式响应、错误处理、测试策略
- **[go-moonshot](https://github.com/northes/go-moonshot)** - Builder 模式、链式调用
- **[deepseek-go](https://github.com/cohesion-org/deepseek-go)** - 环境变量支持、模块化设计

在这些项目的基础上，我们做了以下改进：

1. **NewClientFromEnv()** - 提供专门的环境变量便捷方法
2. **智能重试机制** - 内置自动重试，减少用户代码
3. **优化的 Builder** - 清晰的方法命名（Add vs Set）
4. **完整的错误层次** - 支持精准的错误处理
5. **中文文档** - 完整的中文文档和示例

详细对比分析请查看 [docs/COMPARISON.md](./docs/COMPARISON.md)

## 📄 许可证

MIT License - 查看 [LICENSE](LICENSE) 文件

## 🔗 相关链接

- [360智脑官网](https://ai.360.com)
- [360智脑开放平台](https://ai.360.com/platform)
- [API 文档](https://ai.360.com/platform/docs/overview)
- [问题反馈](https://github.com/lin-coco/zhinao-go/issues)
- [贡献代码](https://github.com/lin-coco/zhinao-go/pulls)

## 💬 支持

如果遇到问题或有建议：

1. 查看 [完整指南](./docs/GUIDE.md) 和 [示例代码](./examples/)
2. 搜索或提交 [Issue](https://github.com/lin-coco/zhinao-go/issues)
3. 参与 [讨论](https://github.com/lin-coco/zhinao-go/discussions)

---

**快速开始**: `go get github.com/lin-coco/zhinao-go`

Made with ❤️ by the community
