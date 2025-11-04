# 360智脑 Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lin-coco/zhinao-go)](https://goreportcard.com/report/github.com/lin-coco/zhinao-go)
[![GoDoc](https://godoc.org/github.com/lin-coco/zhinao-go?status.svg)](https://godoc.org/github.com/lin-coco/zhinao-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

360智脑官方 Go 语言 SDK，提供简洁、类型安全的 API 接口来访问360智脑的 AI 能力。

## ✨ 特性

- 🚀 **简单易用** - 提供直观的 API 接口和链式构建器
- 🔄 **流式支持** - 完整支持流式聊天补全
- ⚡ **高性能** - 内置连接池和智能重试机制
- 🛡️ **类型安全** - 完整的类型定义和编译时检查
- 🔧 **可配置** - 灵活的配置选项，支持自定义 HTTP 客户端
- 📝 **详细文档** - 完善的代码注释和示例
- 🎯 **Context 支持** - 所有 API 都支持 context.Context
- 🔁 **自动重试** - 可配置的指数退避重试策略
- 🌐 **多客户端支持** - 支持标准库和自定义 HTTP 客户端

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
    // 创建客户端
    client, err := zhinao.NewClient("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    // 创建聊天请求
    req := &zhinao.ChatRequest{
        Model: "360gpt-turbo",
        Messages: []zhinao.Message{
            {Role: "user", Content: "你好，请介绍一下360智脑"},
        },
    }

    // 发送请求
    resp, err := client.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    // 打印响应
    fmt.Println(resp.Choices[0].Message.Content)
}
```

### 使用环境变量

SDK 支持通过环境变量 `ZHINAO_API_KEY` 配置 API 密钥：

```bash
# 设置环境变量
export ZHINAO_API_KEY=your-api-key
```

```go
// 方法1: 使用 NewClientFromEnv（推荐）
client, err := zhinao.NewClientFromEnv()

// 方法2: 使用空字符串（自动读取环境变量）
client, err := zhinao.NewClient("")
```

### 使用 Builder 模式

```go
req := zhinao.NewChatBuilder().
    SetModel("360gpt-turbo").
    AddSystemMessage("你是一个专业的编程助手").
    AddUserMessage("什么是 Go 语言的 interface？").
    SetTemperature(0.7).
    SetMaxTokens(500).
    Build()

resp, err := client.Chat.CreateCompletion(ctx, req)
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
    
    if len(resp.Choices) > 0 {
        fmt.Print(resp.Choices[0].Delta.Content)
    }
}
```

### 自定义配置

```go
client, err := zhinao.NewClient(
    "your-api-key",
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithRetry(5, 2*time.Second),
    zhinao.WithBaseURL("https://custom-api.360.cn/v1"),
)
```

## 📖 文档

### 核心概念

#### 客户端配置

SDK 使用函数式选项模式进行配置：

```go
type Option func(*Config) error

// 可用的配置选项：
WithBaseURL(url string)                      // 设置 API 基础 URL
WithTimeout(timeout time.Duration)           // 设置请求超时时间
WithRetry(maxRetries int, delay time.Duration) // 设置重试策略
WithHTTPClient(client http.Client)           // 使用自定义 HTTP 客户端
WithUserAgent(userAgent string)              // 设置 User-Agent
WithHeaders(headers map[string]string)       // 设置自定义请求头
```

#### 聊天服务

聊天服务提供两个主要方法：

- `CreateCompletion` - 非流式聊天补全
- `CreateCompletionStream` - 流式聊天补全

#### 请求构建器

`ChatBuilder` 提供链式 API 来构建复杂的聊天请求：

```go
builder := zhinao.NewChatBuilder()

// 设置模型和参数
builder.SetModel("360gpt-turbo")
builder.SetTemperature(0.7)
builder.SetMaxTokens(1000)

// 添加消息
builder.AddSystemMessage("系统消息")
builder.AddUserMessage("用户消息")
builder.AddAssistantMessage("助手消息")

// 构建请求
req := builder.Build()
```

### 错误处理

SDK 定义了多种错误类型：

```go
// API 错误
type APIError struct {
    StatusCode int
    Type       string
    Code       string
    Message    string
}

// 限流错误
type RateLimitError struct {
    *APIError
    RetryAfter int
}

// 验证错误
type ValidationError struct {
    Field   string
    Message string
}
```

处理错误的示例：

```go
resp, err := client.Chat.CreateCompletion(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *zhinao.APIError:
        fmt.Printf("API错误: %s (状态码: %d)\n", e.Message, e.StatusCode)
    case *zhinao.RateLimitError:
        fmt.Printf("限流错误，请在 %d 秒后重试\n", e.RetryAfter)
    case *zhinao.ValidationError:
        fmt.Printf("验证错误: %s - %s\n", e.Field, e.Message)
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 重试机制

SDK 内置了智能重试机制：

- 支持可配置的最大重试次数
- 使用指数退避策略
- 自动识别可重试的错误（5xx、429、408）
- 遵守 context 取消

```go
// 配置重试：最多重试 5 次，初始延迟 2 秒
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithRetry(5, 2*time.Second),
)
```

### Context 支持

所有 API 都支持 context，可以实现：

- 超时控制
- 取消操作
- 传递请求范围的值

```go
// 设置超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Chat.CreateCompletion(ctx, req)
```

## 🔧 高级用法

### 多轮对话

```go
builder := zhinao.NewChatBuilder().
    SetModel("360gpt-turbo").
    AddSystemMessage("你是一个助手")

// 第一轮
builder.AddUserMessage("什么是 AI？")
resp1, _ := client.Chat.CreateCompletion(ctx, builder.Build())
builder.AddAssistantMessage(resp1.Choices[0].Message.Content)

// 第二轮
builder.AddUserMessage("它有什么用？")
resp2, _ := client.Chat.CreateCompletion(ctx, builder.Build())
```

### 自定义 HTTP 客户端

```go
type MyHTTPClient struct {
    // 你的实现
}

func (c *MyHTTPClient) Post(ctx context.Context, path string, body, result interface{}, apiKey string) error {
    // 实现 http.Client 接口
}

// 使用自定义客户端
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithHTTPClient(&MyHTTPClient{}),
)
```

### 工具调用（Function Calling）

```go
req := &zhinao.ChatRequest{
    Model: "360gpt-turbo",
    Messages: []zhinao.Message{
        {Role: "user", Content: "北京今天天气怎么样？"},
    },
    Tools: []zhinao.Tool{
        {
            Type: "function",
            Function: zhinao.ToolFunction{
                Name:        "get_weather",
                Description: "获取指定城市的天气信息",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "city": map[string]interface{}{
                            "type":        "string",
                            "description": "城市名称",
                        },
                    },
                    "required": []string{"city"},
                },
            },
        },
    },
}
```

## 📚 示例

查看 [examples](./examples) 目录获取更多示例：

- [基础示例](./examples/basic/)
  - 简单聊天
  - Builder 模式使用
- [高级示例](./examples/advanced/)
  - 流式聊天
  - 多轮对话
  - 自定义配置

## 🧪 测试

项目包含完整的单元测试，确保代码质量和稳定性。

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 仅运行单元测试
go test -v .

# 使用 Makefile
make test           # 运行所有测试
make test-unit      # 仅运行单元测试
make test-coverage  # 生成覆盖率报告
```

### 测试覆盖

- ✅ 客户端配置测试
- ✅ 环境变量支持测试
- ✅ Builder 模式测试
- ✅ 错误处理测试
- ✅ 请求验证测试
- ✅ 类型定义测试

### 发布前检查清单

在发布新版本前，请确保：

```bash
# 1. 运行所有测试
make test

# 2. 检查代码风格
make lint

# 3. 生成覆盖率报告
make test-coverage

# 4. 构建项目
make build
```

## 🤝 贡献

欢迎贡献代码！在提交 PR 之前，请确保：

1. ✅ 所有测试通过 (`make test`)
2. ✅ 代码通过 lint 检查 (`make lint`)
3. ✅ 添加必要的测试用例
4. ✅ 更新相关文档
5. ✅ 遵循项目的代码风格

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [360智脑官网](https://ai.360.cn)
- [API 文档](https://ai.360.cn/docs)
- [问题反馈](https://github.com/lin-coco/zhinao-go/issues)

## 📮 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue: https://github.com/lin-coco/zhinao-go/issues

## 🙏 致谢

本项目的设计参考了以下优秀的开源项目：

- [go-moonshot](https://github.com/northes/go-moonshot)
- [deepseek-go](https://github.com/cohesion-org/deepseek-go)
- [go-openai](https://github.com/sashabaranov/go-openai)
