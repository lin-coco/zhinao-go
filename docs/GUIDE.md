# 360智脑 Go SDK - 完整指南

本文档是 360智脑 Go SDK 的完整使用和开发指南，整合了架构设计、最佳实践和测试策略。

## 目录

- [快速开始](#快速开始)
- [架构设计](#架构设计)
- [核心功能](#核心功能)
- [测试指南](#测试指南)
- [最佳实践](#最佳实践)

---

## 快速开始

### 安装

```bash
go get github.com/lin-coco/zhinao-go
```

**要求**: Go >= 1.18

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

---

## 架构设计

### 整体架构

```
┌──────────────────────────────────────────────────────────────┐
│                      客户端层 (Client)                         │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────┐  │
│  │   Chat   │  │  Models  │  │ Embeddings │  │  Images  │  │
│  │  Service │  │  Service │  │   Service  │  │  Service │  │
│  └──────────┘  └──────────┘  └────────────┘  └──────────┘  │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                      HTTP 客户端层                             │
│  ┌──────────────────┐          ┌──────────────────┐          │
│  │  HTTPDoer 接口   │ ◄──────► │  http.Client     │          │
│  │  (可扩展)        │          │  (标准库)        │          │
│  └──────────────────┘          └──────────────────┘          │
│           △                                                    │
│           │                                                    │
│  ┌──────────────────────────────────────────────┐            │
│  │  支持第三方 HTTP 库（通过适配器）               │            │
│  │  - go-resty                                  │            │
│  │  - go-retryablehttp                          │            │
│  │  - 自定义 HTTP 客户端                         │            │
│  └──────────────────────────────────────────────┘            │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                        360智脑 API                             │
└──────────────────────────────────────────────────────────────┘
```

### 项目结构

```
zhinao-go/
├── client.go              # 客户端入口和初始化
├── config.go              # 配置管理（函数式选项模式，包含 HTTPDoer 接口）
├── chat.go               # 聊天补全服务
├── chat_builder.go       # 聊天请求 Builder
├── chat_stream.go        # 流式响应处理
├── models.go             # 模型管理服务
├── embeddings.go         # 向量生成服务
├── images.go             # 图像生成服务
├── types.go              # 公共类型定义
├── constants.go          # 常量定义（模型名、角色等）
├── errors.go             # 错误类型定义
├── internal/             # 内部实现（不导出）
│   └── test/            # 测试工具
│       ├── server.go    # Mock HTTP 服务器
│       └── handlers.go  # 测试请求处理器
├── examples/            # 完整使用示例
│   ├── chat-completion/
│   ├── stream-chat/
│   ├── chat-with-tools/
│   ├── chat-with-builder/
│   ├── chatbot/
│   ├── custom-http-client/  # 自定义 HTTP 客户端示例
│   ├── list-models/
│   ├── embeddings/
│   └── text2img/
└── docs/               # 文档
    ├── README.md       # 文档导航
    ├── GUIDE.md        # 完整指南（本文档）
    └── COMPARISON.md   # SDK 对比分析
```

### 核心设计原则

1. **易用性** - 简洁直观的 API，降低学习成本
2. **类型安全** - 充分利用 Go 的类型系统
3. **可扩展性** - HTTPDoer 接口支持灵活扩展和第三方库集成
4. **高性能** - 优化网络请求和资源使用
5. **可靠性** - 完善的错误处理和超时控制

---

## 核心功能

### 1. 客户端配置

#### 函数式选项模式

```go
client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithTimeout(30*time.Second),
    zhinao.WithBaseURL("https://custom-api.360.cn/v1"),
    zhinao.WithHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    }),
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
        // 根据状态码决定是否重试
        if e.StatusCode >= 500 || e.StatusCode == 429 {
            // 可以考虑重试
        }
    case *zhinao.RateLimitError:
        fmt.Printf("限流错误，请在 %d 秒后重试\n", e.RetryAfter)
        time.Sleep(time.Duration(e.RetryAfter) * time.Second)
    case *zhinao.ValidationError:
        fmt.Printf("验证错误: %s - %s\n", e.Field, e.Message)
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 5. 自定义 HTTP 客户端

SDK 通过 HTTPDoer 接口支持自定义 HTTP 客户端：

```go
// 使用自定义 http.Client
customHTTPClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

client, err := zhinao.NewClient(
    apiKey,
    zhinao.WithHTTPClient(customHTTPClient),
)
```

**支持的场景**：
- 使用标准库 `http.Client`
- 集成第三方 HTTP 库（如 go-resty、go-retryablehttp）
- 实现自定义的 HTTPDoer 接口

详见示例：`examples/custom-http-client/`

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

### 7. 向量生成（Embeddings）

向量生成用于将文本转换为数值向量，支持语义搜索、文本分类、相似度计算等应用。

#### 基础用法

```go
req := &zhinao.EmbeddingsRequest{
    Model: zhinao.ModelEmbeddingS1V1,
    Input: []string{"你好", "世界"},
}

resp, err := client.Embeddings.Create(ctx, req)
if err != nil {
    log.Fatal(err)
}

// 访问向量数据
for i, data := range resp.Data {
    fmt.Printf("文本 %d 的向量维度: %d\n", i, len(data.Embedding))
    fmt.Printf("向量前5维: %v\n", data.Embedding[:5])
}
```

#### 使用 Builder 模式

```go
req := zhinao.NewEmbeddings(zhinao.ModelEmbeddingS1V1).
    AddInput("机器学习").
    AddInput("深度学习").
    AddInput("神经网络").
    SetUser("user-123").
    Build()

resp, err := client.Embeddings.Create(ctx, req)
```

**应用场景**：
- 语义搜索
- 文本分类
- 推荐系统
- 文本聚类
- 问答系统

### 8. 图像生成（Text-to-Image）

文本生成图像功能支持将文本描述转换为图像，提供多种风格选择。

#### 基础用法

```go
req := &zhinao.Text2ImgRequest{
    Model:  zhinao.Model360CVW0V5,
    Style:  zhinao.StyleRealistic,
    Prompt: "一只可爱的小猫在草地上玩耍，阳光明媚",
}

resp, err := client.Images.Text2Img(ctx, req)
if err != nil {
    log.Fatal(err)
}

// 保存图片
for i, img := range resp.Data {
    filename := fmt.Sprintf("image_%d.png", i)
    err := os.WriteFile(filename, img.ImageData, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
```

#### 高级参数配置

```go
req := &zhinao.Text2ImgRequest{
    Model:  zhinao.Model360CVW0V5,
    Style:  zhinao.StyleCartoon,
    Prompt: "一个科幻城市的夜景",
    NegativePrompt: "模糊，低质量",
    Width:  1024,
    Height: 768,
    Samples: 2,
    NumInferenceSteps: 30,
    GuidanceScale: 10.0,
    Seed: 42,
}

resp, err := client.Images.Text2Img(ctx, req)
```

**支持的风格**：
- `StyleRealistic` - 写实风格
- `StyleCartoon` - 卡通风格  
- `StylePapercut` - 剪纸风格
- `StyleCG` - CG 风格

**支持的模型**：
- `Model360CVW0V5` - 360CV W0 V5
- `Model360CVC0V5` - 360CV C0 V5
- `Model360Flux1KontextDev` - Flux 1 Kontext Dev
- 其他第三方模型（详见常量定义）

### 9. 模型管理

SDK 提供模型列表查询和详情获取功能。

#### 获取模型列表

```go
models, err := client.Models.List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, model := range models.Data {
    fmt.Printf("模型ID: %s\n", model.ID)
    fmt.Printf("  名称: %s\n", model.Name)
    fmt.Printf("  拥有者: %s\n", model.OwnedBy)
    fmt.Printf("  类型: %s\n", model.Object)
    fmt.Println()
}
```

#### 获取模型详情

```go
model, err := client.Models.Get(ctx, zhinao.Model360GPTTurbo)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("模型: %s\n", model.ID)
fmt.Printf("创建时间: %d\n", model.Created)
```

**常用模型**：
- 聊天模型：`Model360GPTTurbo`, `Model360GPTPro`, `ModelDeepSeekV3` 等
- 图像模型：`Model360CVW0V5`, `ModelDallE3` 等
- 向量模型：`ModelEmbeddingS1V1`

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

### 5. 错误处理

SDK 提供了完善的错误类型供精准处理：

```go
resp, err := client.Chat.CreateCompletion(ctx, req)
if err != nil {
    switch e := err.(type) {
    case *zhinao.APIError:
        log.Printf("API错误: %s (状态码: %d)\n", e.Message, e.StatusCode)
    case *zhinao.RateLimitError:
        log.Printf("限流错误，建议等待 %d 秒后重试\n", e.RetryAfter)
    case *zhinao.ValidationError:
        log.Printf("验证错误: %s - %s\n", e.Field, e.Message)
    default:
        log.Printf("未知错误: %v\n", err)
    }
    return
}
```

**错误类型说明**：
- `APIError` - API 返回的错误，包含状态码和详细信息
- `RateLimitError` - 限流错误，包含建议的等待时间
- `ValidationError` - 请求验证错误
- 预定义错误 - `ErrMissingAPIKey`、`ErrInvalidModel` 等

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
