# SDK 设计对比分析

本文档对比 360智脑 Go SDK 与其他流行 Go SDK 的设计差异，帮助开发者理解不同的设计选择。

## 📋 对比项目

- **[go-openai](https://github.com/sashabaranov/go-openai)** - OpenAI 官方 Go SDK
- **[go-moonshot](https://github.com/northes/go-moonshot)** - Moonshot AI (Kimi) SDK
- **[deepseek-go](https://github.com/cohesion-org/deepseek-go)** - DeepSeek SDK

## 🔧 客户端配置对比

### zhinao-go

```go
// 方式1: 环境变量便捷方法
client, err := zhinao.NewClientFromEnv()

// 方式2: 自动回退到环境变量
client, err := zhinao.NewClient("")

// 方式3: 显式传入 API Key
client, err := zhinao.NewClient("your-api-key")

// 方式4: 带配置选项
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithBaseURL("https://api.360.cn/v1"),
)
```

**特点**：
- 提供专门的 `NewClientFromEnv()` 方法
- 空字符串自动回退到环境变量 `ZHINAO_API_KEY`
- 函数式选项模式，向后兼容

### go-moonshot

```go
client, err := moonshot.NewClient(apiKey)
```

**特点**：
- 不支持环境变量
- 需要用户手动调用 `os.Getenv()`

### deepseek-go

```go
// 空字符串时自动读取环境变量
client := deepseek.NewClient("")

// 使用选项
client, err := deepseek.NewClientWithOptions("",
    deepseek.WithTimeout(5*time.Minute),
)
```

**特点**：
- 支持环境变量 `DEEPSEEK_API_KEY`（使用 `os.LookupEnv`）
- 空字符串自动回退
- 返回 nil 而不是 error

### go-openai

```go
client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

config := openai.DefaultConfig(apiKey)
client := openai.NewClientWithConfig(config)
```

**特点**：
- 需要用户显式调用 `os.Getenv()`
- 配置结构体模式

## 🏗️ 请求构建对比

### zhinao-go - Builder 模式

```go
req := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddSystemMessage("你是一个助手").
    AddUserMessage("你好").
    SetTemperature(0.7).
    SetMaxTokens(1000).
    Build()
```

**特点**：
- 链式调用
- 方法命名清晰（Add vs Set）
- 易于维护对话历史

### go-moonshot - Builder 支持

```go
builder := moonshot.NewChatCompletionsBuilder()
builder.AppendPrompt("系统提示").
    AppendUser("用户消息").
    WithTemperature(0.3)

req := builder.ToRequest()
```

**特点**：
- 支持 Builder
- 使用 Append 系列方法

### deepseek-go - 直接构造

```go
req := &deepseek.ChatCompletionRequest{
    Model: deepseek.DeepSeekChat,
    Messages: []deepseek.ChatCompletionMessage{
        {Role: deepseek.ChatMessageRoleSystem, Content: "系统"},
        {Role: deepseek.ChatMessageRoleUser, Content: "用户"},
    },
    Temperature: 0.7,
}
```

**特点**：
- 无 Builder 模式
- 直接构造结构体

### go-openai - 直接构造

```go
req := openai.ChatCompletionRequest{
    Model: openai.GPT4,
    Messages: []openai.ChatCompletionMessage{
        {Role: openai.ChatMessageRoleSystem, Content: "系统"},
        {Role: openai.ChatMessageRoleUser, Content: "用户"},
    },
}
```

**特点**：
- 无 Builder 模式
- 结构体直接构造

## 🌊 流式响应对比

所有 SDK 都采用标准的迭代器模式：

```go
stream, err := client.CreateCompletionStream(ctx, req)
defer stream.Close()

for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        return err
    }
    // 处理响应
}
```

这是 Go 社区的通用做法。

## ⚠️ 错误处理对比

### zhinao-go

```go
type APIError struct {
    StatusCode int
    Type       string
    Code       string
    Message    string
}

type RateLimitError struct {
    *APIError
    RetryAfter int
}

type ValidationError struct {
    Field   string
    Message string
}

// 预定义错误
var (
    ErrMissingAPIKey = errors.New("missing API key")
    ErrInvalidModel  = errors.New("invalid model")
    // ...
)
```

**使用示例**：
```go
switch e := err.(type) {
case *zhinao.APIError:
    log.Printf("API错误: %s", e.Error())
case *zhinao.RateLimitError:
    time.Sleep(time.Duration(e.RetryAfter) * time.Second)
case *zhinao.ValidationError:
    log.Printf("字段 %s 验证失败: %s", e.Field, e.Message)
}
```

**特点**：
- 错误类型层次结构
- 预定义常见错误
- 提供详细的错误信息

### go-moonshot

```go
type APIError struct {
    StatusCode int
    Message    string
    Type       string
}
```

**特点**：
- 基础错误结构

### deepseek-go

```go
type Error struct {
    Code    string
    Message string
}

func HandleAPIError(resp *http.Response) error
```

**特点**：
- 简单的错误结构
- 提供错误处理函数

### go-openai

```go
type APIError struct {
    StatusCode int
    Message    string
    Type       string
    Code       string
}

type RateLimitError struct {
    *APIError
}
```

**特点**：
- 包含限流错误类型
- 错误信息较完整

## 📊 功能对比

| 功能 | zhinao-go | go-moonshot | deepseek-go | go-openai |
|-----|-----------|-------------|-------------|-----------|
| 聊天补全 | ✅ | ✅ | ✅ | ✅ |
| 流式响应 | ✅ | ✅ | ✅ | ✅ |
| 工具调用 | ✅ | ✅ | ✅ | ✅ |
| 模型列表 | ✅ | ✅ | ✅ | ✅ |
| 向量生成 | ✅ | ❌ | ❌ | ✅ |
| 图像生成 | ✅ | ❌ | ✅ | ✅ |
| 环境变量便捷方法 | ✅ | ❌ | ❌ | ❌ |
| Builder 模式 | ✅ | ✅ | ❌ | ❌ |
| HTTPDoer 接口 | ✅ | ❌ | ❌ | ✅ |
| 中文文档 | ✅ | ✅ | ❌ | ❌ |

## 🎨 设计借鉴

### 从 go-openai 学习

- 流式响应的标准迭代器模式
- 完整的类型定义
- Mock Server 测试策略

### 从 go-moonshot 学习

- Builder 模式的应用
- 链式调用的流畅性
- 清晰的 API 命名

### 从 deepseek-go 学习

- 环境变量的处理方式（使用 `os.LookupEnv`）
- 模块化的代码组织
- HTTP 客户端抽象，轻松集成第三方 HTTP 库

## 🆕 zhinao-go 的特点

基于以上学习，zhinao-go 做了以下设计：

1. **NewClientFromEnv()** - 提供专门的环境变量便捷方法
2. **HTTPDoer 接口** - 灵活的 HTTP 客户端抽象，可轻松集成第三方 HTTP 库（如 go-resty、go-retryablehttp）
3. **优化的 Builder** - 清晰的方法命名（Add vs Set）
4. **完整的错误层次** - 支持精准的错误处理
5. **中文文档** - 完整的中文文档和示例

## 📚 参考资料

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

如有问题或建议，欢迎通过 [GitHub Issues](https://github.com/lin-coco/zhinao-go/issues) 反馈。
