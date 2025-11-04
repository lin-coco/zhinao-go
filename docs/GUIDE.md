# 360智脑 Go SDK - 完整指南

本文档是 360智脑 Go SDK 的完整使用和开发指南，整合了架构设计、最佳实践和测试策略。

## 目录

- [快速开始](#快速开始)
- [架构设计](#架构设计)
- [核心功能](#核心功能)
- [测试指南](#测试指南)
- [最佳实践](#最佳实践)
- [扩展开发](#扩展开发)

---

## 快速开始

### 安装

```bash
go get github.com/lin-coco/zhinao-go
```

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
    // 创建客户端（自动从环境变量读取 ZHINAO_API_KEY）
    client, err := zhinao.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // 使用 Builder 构建请求
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

---

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    客户端层 (Client)                      │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │ ChatService │  │ ModelsService│  │ Future Services│  │
│  └─────────────┘  └──────────────┘  └───────────────┘  │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                 HTTP 客户端层 (internal/http)            │
│  ┌──────────────────┐         ┌───────────────────┐    │
│  │ StandardClient   │ ◄─────► │ Client Interface  │    │
│  └──────────────────┘         └───────────────────┘    │
│            │                           △                │
│            ▼                           │                │
│  ┌──────────────────┐         ┌───────────────────┐    │
│  │  Retry Logic     │         │  Custom Clients   │    │
│  └──────────────────┘         └───────────────────┘    │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                     360智脑 API                          │
└─────────────────────────────────────────────────────────┘
```

### 项目结构

```
zhinao-go/
├── client.go              # 客户端入口和配置
├── config.go              # 配置管理（函数式选项）
├── chat.go               # 聊天服务实现
├── chat_builder.go       # 链式构建器
├── chat_stream.go        # 流式响应处理
├── models.go             # 模型管理
├── types.go              # 类型定义
├── constants.go          # 常量定义
├── errors.go             # 错误类型
├── internal/             # 内部实现
│   ├── http/            # HTTP 客户端抽象
│   │   ├── client.go    # 接口定义
│   │   ├── standard.go  # 标准实现
│   │   └── retry.go     # 重试逻辑
│   └── test/            # 测试工具
│       ├── server.go    # Mock 服务器
│       └── handlers.go  # 请求处理器
├── examples/            # 使用示例
└── docs/               # 文档
```

### 核心设计原则

1. **易用性** - 简洁直观的 API，降低学习成本
2. **类型安全** - 充分利用 Go 的类型系统
3. **可扩展性** - 支持未来功能扩展
4. **高性能** - 优化网络请求和资源使用
5. **可靠性** - 内置重试机制和错误处理

---

## 核心功能

### 1. 客户端配置

#### 函数式选项模式

```go
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithRetry(5, 2*time.Second),
    zhinao.WithBaseURL("https://custom-api.360.cn/v1"),
    zhinao.WithUserAgent("MyApp/1.0"),
)
```

**优点**：
- 向后兼容：添加新选项不会破坏现有代码
- 可选参数：只需设置需要的参数
- 默认值：提供合理的默认配置
- 链式调用：提供良好的用户体验

#### 环境变量支持

```go
// 方式1: 使用专门的便捷方法（推荐）
client, err := zhinao.NewClientFromEnv()

// 方式2: 传入空字符串，自动读取环境变量
client, err := zhinao.NewClient("")

// 方式3: 结合配置选项
client, err := zhinao.NewClientFromEnv(
    zhinao.WithTimeout(30*time.Second),
)
```

环境变量：`ZHINAO_API_KEY`

### 2. Builder 模式

Builder 模式简化复杂请求的构建：

```go
req := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddSystemMessage("你是一个专业的技术顾问").
    AddUserMessage("什么是机器学习？").
    SetTemperature(0.7).
    SetMaxTokens(1000).
    SetTopP(0.9).
    Build()
```

**优势**：
- 流畅的 API
- 参数验证
- 不可变性
- 可读性强
- 便于维护对话历史

### 3. 流式响应

```go
req := &zhinao.ChatRequest{
    Model:  zhinao.Model360GPTTurbo,
    Messages: []zhinao.Message{
        {Role: zhinao.RoleUser, Content: "写一首诗"},
    },
    Stream: true,
}

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

### 4. 错误处理

#### 错误类型层次

```
error (interface)
  │
  ├── APIError (API 错误)
  │     └── RateLimitError (限流错误)
  │
  ├── ValidationError (验证错误)
  │
  └── 预定义错误
       ├── ErrMissingAPIKey
       ├── ErrInvalidModel
       ├── ErrEmptyMessages
       └── ErrStreamClosed
```

#### 错误处理示例

```go
resp, err := client.Chat.CreateCompletion(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *zhinao.APIError:
        fmt.Printf("API错误: %s (状态码: %d)\n", e.Message, e.StatusCode)
        if e.IsRetryable() {
            // 可以重试
        }
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

### 5. 自动重试机制

SDK 内置智能重试机制：

```go
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithRetry(5, 2*time.Second), // 最多重试5次，初始延迟2秒
)
```

**特性**：
- 自动识别可重试错误（5xx、429、408）
- 指数退避策略
- 遵守 context 取消
- 遵守服务器的 Retry-After 头

### 6. 工具调用（Function Calling）

```go
req := &zhinao.ChatRequest{
    Model: zhinao.Model360GPTTurbo,
    Messages: []zhinao.Message{
        {Role: zhinao.RoleUser, Content: "北京今天天气怎么样？"},
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

// 第一次调用：AI 决定是否使用工具
resp, err := client.Chat.CreateCompletion(ctx, req)

// 检查工具调用
if len(resp.Choices[0].Message.ToolCalls) > 0 {
    // 执行工具调用，然后将结果返回给 AI
    // ...（详见示例 examples/chat-with-tools）
}
```

---

## 测试指南

### 测试策略

项目采用分层测试策略：

```
单元测试 ✅
├── 客户端配置测试
├── Builder 模式测试
├── 错误处理测试
└── 类型验证测试

Mock Server 测试 ✅
├── 聊天补全测试
├── 流式响应测试
├── 错误响应测试
├── 限流处理测试
└── 超时重试测试

集成测试（可选）
└── 真实 API 调用
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行单元测试
make test-unit

# 生成覆盖率报告
make test-coverage

# 查看覆盖率
make test-coverage
# 然后在浏览器打开 coverage.html
```

### Mock Server 测试

SDK 使用 `httptest.Server` 创建模拟服务器，**不需要真实 API Key**：

```go
func TestChatCompletion_Mock(t *testing.T) {
    // 创建 Mock 服务器
    server := internal_test.NewMockServer()
    defer server.Close()
    
    // 注册响应处理器
    server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
        resp := &zhinao.ChatResponse{
            ID: "test-id",
            Choices: []zhinao.Choice{
                {
                    Message: zhinao.Message{
                        Role:    zhinao.RoleAssistant,
                        Content: "Hello!",
                    },
                },
            },
        }
        json.NewEncoder(w).Encode(resp)
    })
    
    // 创建指向 Mock 服务器的客户端
    client, _ := zhinao.NewClient(
        "test-api-key",
        zhinao.WithBaseURL(server.URL),
    )
    
    // 执行测试
    req := &zhinao.ChatRequest{
        Model: zhinao.Model360GPTTurbo,
        Messages: []zhinao.Message{
            {Role: zhinao.RoleUser, Content: "Hi"},
        },
    }
    
    resp, err := client.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if resp.Choices[0].Message.Content != "Hello!" {
        t.Errorf("Unexpected response")
    }
}
```

### Mock Server 的优势

- ✅ **无需 API Key** - 测试可在任何环境运行
- ✅ **快速可靠** - 不依赖网络，秒级完成
- ✅ **可预测性** - 完全控制响应内容
- ✅ **边界测试** - 可模拟各种错误情况
- ✅ **成本为零** - 不消耗 API 配额
- ✅ **CI/CD 友好** - 可在 GitHub Actions 等环境运行

### 测试最佳实践

1. **测试独立性** - 每个测试应该独立，不依赖其他测试

2. **表格驱动测试** - 对于多个类似场景：
```go
tests := []struct {
    name    string
    input   string
    wantErr bool
}{
    {"valid", "test", false},
    {"empty", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // 测试逻辑
    })
}
```

3. **使用子测试** - 组织相关测试用例：
```go
func TestClient(t *testing.T) {
    t.Run("with valid key", func(t *testing.T) {
        // ...
    })
    t.Run("with empty key", func(t *testing.T) {
        // ...
    })
}
```

4. **清理资源** - 使用 defer 确保资源清理：
```go
server := setupTestServer()
defer server.Close()
```

---

## 最佳实践

### 1. 使用环境变量管理 API Key

```bash
# 设置环境变量
export ZHINAO_API_KEY="your-api-key"
```

```go
// 代码中使用
client, err := zhinao.NewClientFromEnv()
```

**优势**：
- 安全：不在代码中硬编码密钥
- 灵活：不同环境使用不同密钥
- 标准：符合 12-Factor App 原则

### 2. 使用 Builder 模式构建请求

```go
// 推荐：使用 Builder
req := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddUserMessage("你好").
    Build()

// 不推荐：直接构造
req := &zhinao.ChatRequest{
    Model: "360gpt-turbo",
    Messages: []zhinao.Message{
        {Role: "user", Content: "你好"},
    },
}
```

### 3. 正确处理流式响应

```go
stream, err := client.Chat.CreateCompletionStream(ctx, req)
if err != nil {
    return err
}
defer stream.Close() // 重要：确保关闭

for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break // 正常结束
    }
    if err != nil {
        return err // 错误处理
    }
    // 处理响应
}
```

### 4. 使用 Context 控制超时

```go
// 设置超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.Chat.CreateCompletion(ctx, req)
```

### 5. 合理配置重试策略

```go
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithRetry(3, 1*time.Second), // 生产环境：适度重试
)
```

**建议**：
- 开发环境：少重试或不重试
- 生产环境：3-5 次重试
- 关键业务：更多重试 + 告警

### 6. 维护对话历史

使用 Builder 模式便于维护多轮对话：

```go
builder := zhinao.NewChatBuilder().
    SetModel(zhinao.Model360GPTTurbo).
    AddSystemMessage("你是一个助手")

// 第一轮
builder.AddUserMessage("什么是AI？")
resp1, _ := client.Chat.CreateCompletion(ctx, builder.Build())
builder.AddAssistantMessage(resp1.Choices[0].Message.Content)

// 第二轮
builder.AddUserMessage("它有什么用？")
resp2, _ := client.Chat.CreateCompletion(ctx, builder.Build())
```

---

## 扩展开发

### 添加新服务

SDK 架构支持便捷地添加新服务（如向量、图像生成）：

```go
// 1. 定义服务接口
type EmbeddingsService interface {
    Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}

// 2. 实现服务
type embeddingsService struct {
    httpClient http.Client
}

func (s *embeddingsService) Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
    // 实现逻辑
}

// 3. 在 Client 中添加
type Client struct {
    Chat       ChatService
    Models     ModelsService
    Embeddings EmbeddingsService // 新服务
}

// 4. 在 NewClient 中初始化
func NewClient(apiKey string, opts ...Option) (*Client, error) {
    // ...
    client.Embeddings = &embeddingsService{httpClient: httpClient}
    return client, nil
}
```

### 自定义 HTTP 客户端

```go
type MyHTTPClient struct {
    // 自定义实现
}

func (c *MyHTTPClient) Post(ctx context.Context, path string, body, result interface{}, apiKey string) error {
    // 实现 http.Client 接口
}

client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithHTTPClient(&MyHTTPClient{}),
)
```

### 中间件支持

```go
type Middleware func(http.Client) http.Client

func LoggingMiddleware(next http.Client) http.Client {
    return &loggingClient{next: next}
}
```

---

## 参考项目

本 SDK 的设计参考了以下优秀项目：

- **go-openai** - 流式响应、配置管理、测试策略
- **go-moonshot** - Builder 模式、错误处理
- **deepseek-go** - 项目结构、接口设计

详细对比分析请参考 [COMPARISON.md](./COMPARISON.md)

---

## 贡献指南

欢迎贡献！在添加新功能时，请遵循：

1. 保持接口简洁
2. 提供完整的文档
3. 编写测试用例
4. 遵循现有代码风格
5. 考虑向后兼容性

### 提交前检查清单

- [ ] 代码通过 `make test`
- [ ] 代码通过 `make lint`
- [ ] 添加了测试用例
- [ ] 更新了相关文档
- [ ] 提交信息清晰明了

---

## 总结

360智脑 Go SDK 提供：

- ✅ 简洁直观的 API
- ✅ 完善的错误处理
- ✅ 内置重试机制
- ✅ 流式响应支持
- ✅ Builder 模式
- ✅ 环境变量支持
- ✅ 完整的测试覆盖
- ✅ 详细的文档和示例
- ✅ 良好的扩展性

这使得 SDK 既易于使用，又便于维护和扩展，适合从快速原型到生产环境的各种场景。
