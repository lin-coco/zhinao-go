# 测试策略详解

## go-openai 的测试方法分析

经过对行业标杆项目 [go-openai](https://github.com/sashabaranov/go-openai) 的深入研究，我们发现其测试策略主要包含以下几个层次：

### 1. **Mock HTTP Server 测试** ⭐ 核心策略

go-openai 使用 `httptest.Server` 创建模拟服务器，**不需要真实调用 API**：

```go
// 测试服务器设置
func setupOpenAITestServer() (client *openai.Client, server *test.ServerTest, teardown func()) {
    server = test.NewTestServer()
    ts := server.OpenAITestServer()
    ts.Start()
    teardown = ts.Close
    
    config := openai.DefaultConfig(test.GetTestToken())
    config.BaseURL = ts.URL + "/v1"  // 指向测试服务器
    client = openai.NewClientWithConfig(config)
    return
}

// 测试示例
func TestChatCompletions(t *testing.T) {
    client, server, teardown := setupOpenAITestServer()
    defer teardown()
    
    // 注册模拟的响应处理器
    server.RegisterHandler("/v1/chat/completions", handleChatCompletionEndpoint)
    
    // 发送请求（实际请求到测试服务器）
    _, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleUser, Content: "Hello!"},
        },
    })
    
    checks.NoError(t, err, "CreateChatCompletion error")
}
```

### 2. **测试服务器实现**

go-openai 在 `internal/test/server.go` 中实现了完整的测试服务器：

```go
type ServerTest struct {
    handlers map[string]handler
}

func (ts *ServerTest) OpenAITestServer() *httptest.Server {
    return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 验证认证
        if r.Header.Get("Authorization") != "Bearer "+GetTestToken() {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        // 路由匹配和处理
        for route, handler := range ts.handlers {
            pattern, _ := regexp.Compile("^" + route + "$")
            if pattern.MatchString(r.URL.Path) {
                handler(w, r)
                return
            }
        }
        
        http.Error(w, "the resource path doesn't exist", http.StatusNotFound)
    }))
}
```

### 3. **测试层次结构**

```
单元测试（已完成）
├── 客户端配置
├── Builder 模式
├── 错误处理
└── 类型验证

Mock Server 测试（推荐添加）✨
├── 聊天补全测试
├── 流式响应测试
├── 错误响应测试
├── 限流处理测试
└── 超时重试测试

集成测试（可选）
└── 真实 API 调用（需要 API Key）
```

## 为什么 Mock Server 测试优于真实 API 测试

### ✅ Mock Server 的优势

1. **无需 API Key** - 测试可以在任何环境运行
2. **快速可靠** - 不依赖网络，秒级完成
3. **可预测性** - 完全控制响应内容
4. **边界测试** - 可以模拟各种错误情况
5. **成本为零** - 不消耗 API 配额
6. **CI/CD 友好** - 可以在 GitHub Actions 等环境运行

### ⚠️ 真实 API 测试的局限

1. **需要密钥** - 增加配置复杂度
2. **不稳定** - 依赖网络和服务可用性  
3. **缓慢** - 网络延迟影响测试速度
4. **成本** - 消耗 API 配额或产生费用
5. **难以测试错误** - 很难触发特定错误情况
6. **环境依赖** - CI/CD 需要额外配置

## 推荐的测试策略

### 阶段一：单元测试（✅ 已完成）

```bash
# 当前已实现
- client_test.go        # 客户端管理
- chat_builder_test.go  # Builder 模式
- errors_test.go        # 错误处理
- types_test.go         # 类型验证
```

### 阶段二：Mock Server 测试（🎯 推荐实现）

```bash
# 推荐添加
- chat_test.go          # 聊天补全 Mock 测试
- chat_stream_test.go   # 流式响应 Mock 测试
- models_test.go        # 模型列表 Mock 测试
- internal/test/        # 测试服务器实现
  ├── server.go         # Mock Server
  └── handlers.go       # 响应处理器
```

### 阶段三：集成测试（📝 可选）

```bash
# 可选实现（需要真实 API Key）
- integration/
  └── api_test.go       # 真实 API 集成测试
```

## Mock Server 实现示例

### 1. 创建测试服务器

```go
// internal/test/server.go
package test

import (
    "net/http"
    "net/http/httptest"
    "regexp"
)

type ServerTest struct {
    handlers map[string]http.HandlerFunc
}

func NewTestServer() *ServerTest {
    return &ServerTest{
        handlers: make(map[string]http.HandlerFunc),
    }
}

func (s *ServerTest) RegisterHandler(path string, handler http.HandlerFunc) {
    s.handlers[path] = handler
}

func (s *ServerTest) Start() *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 验证 API Key
        if r.Header.Get("Authorization") != "Bearer test-api-key" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        // 路由匹配
        for path, handler := range s.handlers {
            matched, _ := regexp.MatchString(path, r.URL.Path)
            if matched {
                handler(w, r)
                return
            }
        }
        
        http.Error(w, "Not Found", http.StatusNotFound)
    }))
}
```

### 2. 聊天补全测试

```go
// chat_test.go
package zhinao

import (
    "context"
    "encoding/json"
    "net/http"
    "testing"
)

func TestChatCompletion(t *testing.T) {
    // 设置测试服务器
    server := setupTestServer()
    defer server.Close()
    
    // 注册处理器
    server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
        // 模拟响应
        resp := ChatResponse{
            ID: "chatcmpl-123",
            Choices: []Choice{
                {
                    Index: 0,
                    Message: Message{
                        Role: "assistant",
                        Content: "Hello! How can I help you?",
                    },
                    FinishReason: "stop",
                },
            },
        }
        json.NewEncoder(w).Encode(resp)
    })
    
    // 创建客户端（指向测试服务器）
    client, _ := NewClient("test-api-key", WithBaseURL(server.URL))
    
    // 执行测试
    req := &ChatRequest{
        Model: "360gpt-turbo",
        Messages: []Message{
            {Role: "user", Content: "Hello"},
        },
    }
    
    resp, err := client.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        t.Fatalf("CreateCompletion failed: %v", err)
    }
    
    if resp.Choices[0].Message.Content != "Hello! How can I help you?" {
        t.Errorf("Unexpected response: %s", resp.Choices[0].Message.Content)
    }
}
```

### 3. 错误处理测试

```go
func TestChatCompletion_RateLimit(t *testing.T) {
    server := setupTestServer()
    defer server.Close()
    
    // 模拟限流错误
    server.RegisterHandler("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Retry-After", "60")
        w.WriteHeader(http.StatusTooManyRequests)
        json.NewEncoder(w).Encode(ErrorResponse{
            Error: ErrorDetail{
                Message: "Rate limit exceeded",
                Type: "rate_limit_error",
                Code: "rate_limit",
            },
        })
    })
    
    client, _ := NewClient("test-api-key", WithBaseURL(server.URL))
    
    _, err := client.Chat.CreateCompletion(context.Background(), &ChatRequest{
        Model: "360gpt-turbo",
        Messages: []Message{{Role: "user", Content: "test"}},
    })
    
    // 验证返回正确的错误类型
    var rateLimitErr *RateLimitError
    if !errors.As(err, &rateLimitErr) {
        t.Error("Expected RateLimitError")
    }
    
    if rateLimitErr.RetryAfter != 60 {
        t.Errorf("Expected RetryAfter=60, got %d", rateLimitErr.RetryAfter)
    }
}
```

## 实施建议

### 当前状态 ✅

你的项目已经有了良好的单元测试基础：
- 客户端管理测试 ✅
- Builder 模式测试 ✅
- 错误处理测试 ✅
- 类型验证测试 ✅

### 下一步行动 🎯

**优先级 1（强烈推荐）：**
1. 实现 Mock Server 基础设施（`internal/test/server.go`）
2. 添加聊天补全的 Mock 测试
3. 添加流式响应的 Mock 测试
4. 添加错误场景的 Mock 测试

**优先级 2（可选）：**
1. 添加集成测试框架（需要真实 API Key）
2. 添加性能基准测试（Benchmark）

### 实施步骤

```bash
# 步骤 1: 创建测试服务器
mkdir -p internal/test
# 实现 internal/test/server.go

# 步骤 2: 添加 Mock 测试
# 实现 chat_test.go（Mock 版本）
# 实现 chat_stream_test.go（Mock 版本）

# 步骤 3: 运行测试
make test

# 步骤 4: 检查覆盖率
make test-coverage
```

## 总结

**go-openai 的成功经验告诉我们：**

1. ✅ **Mock Server 是主要测试方式** - 不依赖真实 API
2. ✅ **快速、可靠、可重复** - 完美适合 CI/CD
3. ✅ **完整的错误场景测试** - 可以模拟各种边界情况
4. ⚠️ **真实 API 测试是补充** - 仅用于最终验证

**你的项目现状：**
- ✅ 单元测试已经很完善
- 🎯 建议添加 Mock Server 测试层
- 📝 真实 API 测试可作为可选项

这种分层测试策略既保证了代码质量，又保持了测试的快速和可靠性！
