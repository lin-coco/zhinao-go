# 360智脑 Go SDK 架构设计

## 概述

本文档详细说明了 360智脑 Go SDK 的架构设计、设计原则和实现细节。

## 设计目标

1. **易用性** - 提供简洁直观的 API，降低学习成本
2. **类型安全** - 充分利用 Go 的类型系统，在编译时捕获错误
3. **可扩展性** - 支持未来功能扩展（向量、图像生成等）
4. **高性能** - 优化网络请求和资源使用
5. **可靠性** - 内置重试机制和错误处理

## 架构图

```
┌─────────────────────────────────────────────────────────┐
│                      客户端层 (Client)                    │
│  ┌─────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │ ChatService │  │ ModelsService│  │ Future Services│  │
│  └─────────────┘  └──────────────┘  └───────────────┘  │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                   HTTP 客户端层 (internal/http)          │
│  ┌──────────────────┐         ┌───────────────────┐    │
│  │ StandardClient   │ ◄─────► │ Client Interface  │    │
│  └──────────────────┘         └───────────────────┘    │
│            │                           △                │
│            │                           │                │
│            ▼                           │                │
│  ┌──────────────────┐         ┌───────────────────┐    │
│  │  Retry Logic     │         │  Custom Clients   │    │
│  └──────────────────┘         └───────────────────┘    │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────┐
│                   360智脑 API                            │
└─────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. 客户端 (Client)

**文件**: `client.go`

客户端是 SDK 的入口点，负责：
- 配置管理
- 服务实例化
- HTTP 客户端初始化

```go
type Client struct {
    config     *Config
    httpClient http.Client
    Chat       ChatService
    Models     ModelsService
}
```

### 2. 配置系统 (Config)

**文件**: `config.go`

采用**函数式选项模式 (Functional Options Pattern)**：

优点：
- 向后兼容：添加新选项不会破坏现有代码
- 可选参数：只需设置需要的参数
- 默认值：提供合理的默认配置
- 链式调用：提供良好的用户体验

```go
type Option func(*Config) error

func WithTimeout(timeout time.Duration) Option {
    return func(c *Config) error {
        c.Timeout = timeout
        return nil
    }
}
```

### 3. 类型系统 (Types)

**文件**: `types.go`

定义了所有的请求和响应类型：
- 使用 struct tags 支持 JSON 序列化
- 提供验证方法确保数据完整性
- 清晰的文档注释

### 4. HTTP 客户端层

**文件**: `internal/http/`

#### 接口设计

```go
type Client interface {
    Post(ctx context.Context, path string, body, result interface{}, apiKey string) error
    Get(ctx context.Context, path string, result interface{}, apiKey string) error
    PostStream(ctx context.Context, path string, body interface{}, apiKey string) (StreamReader, error)
}
```

**设计考虑**：
- 接口抽象允许多种实现（标准库、第三方库）
- Context 支持超时和取消
- 流式接口独立设计

#### 重试机制

**文件**: `internal/http/retry.go`

特性：
- 可配置的最大重试次数
- 指数退避策略
- 智能判断可重试的错误
- 遵守 context 取消

```go
type RetryConfig struct {
    MaxRetries int
    RetryDelay time.Duration
    Backoff    BackoffStrategy
}
```

### 5. 服务层

#### Chat 服务

**文件**: `chat.go`, `chat_stream.go`, `chat_builder.go`

**设计模式**：

1. **服务接口模式** - 定义清晰的服务边界
2. **Builder 模式** - 简化复杂请求的构建
3. **Iterator 模式** - 流式响应使用迭代器模式

```go
type ChatService interface {
    CreateCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    CreateCompletionStream(ctx context.Context, req *ChatRequest) (ChatStream, error)
}
```

#### Builder 模式优势

```go
req := NewChatBuilder().
    SetModel("360gpt-turbo").
    AddUserMessage("Hello").
    SetTemperature(0.7).
    Build()
```

- 流畅的 API
- 参数验证
- 不可变性（返回新实例）
- 可读性强

### 6. 错误处理

**文件**: `errors.go`

**错误层次结构**：

```
error (interface)
  │
  ├── APIError (API 错误)
  │     └── RateLimitError (限流错误)
  │
  ├── ValidationError (验证错误)
  │
  └── 其他标准错误
```

**特性**：
- 类型断言支持
- 错误上下文信息
- 可重试判断方法

## 设计模式

### 1. 函数式选项模式 (Functional Options)

**应用**: 客户端配置

**优点**:
- API 稳定性
- 向后兼容
- 清晰的意图

### 2. Builder 模式

**应用**: 请求构建

**优点**:
- 链式调用
- 参数验证
- 灵活性

### 3. 接口分离原则 (Interface Segregation)

**应用**: HTTP 客户端、服务接口

**优点**:
- 依赖注入
- 可测试性
- 多实现支持

### 4. 依赖注入 (Dependency Injection)

**应用**: HTTP 客户端注入到服务

**优点**:
- 松耦合
- 易于测试
- 灵活替换实现

## 扩展性设计

### 未来功能支持

SDK 的架构设计考虑了未来功能扩展：

```go
type Client struct {
    config     *Config
    httpClient http.Client
    
    // 当前服务
    Chat   ChatService
    Models ModelsService
    
    // 未来服务（示例）
    // Embeddings EmbeddingsService  // 向量服务
    // Images     ImagesService       // 图像生成服务
    // Audio      AudioService        // 音频服务
}
```

添加新服务的步骤：

1. 定义服务接口
2. 实现服务接口
3. 在 Client 中添加服务字段
4. 在 NewClient 中初始化服务

### 中间件支持

HTTP 客户端设计支持中间件扩展：

```go
type Middleware func(Client) Client

// 示例：日志中间件
func LoggingMiddleware(next Client) Client {
    return &loggingClient{next: next}
}
```

## 性能优化

### 1. 连接复用

标准 HTTP 客户端自动复用连接：

```go
httpClient: &http.Client{
    Timeout: timeout,
    // 默认会复用连接
}
```

### 2. Context 使用

所有 API 都支持 context：
- 超时控制
- 及时取消
- 资源释放

### 3. 流式响应

流式 API 设计：
- 增量处理
- 减少内存占用
- 提高响应速度

## 测试策略

### 单元测试

- Mock HTTP 客户端
- 测试各个组件独立功能
- 测试错误处理

### 集成测试

- 真实 API 调用（可选）
- 端到端测试
- 性能测试

### 示例代码

`examples/` 目录提供：
- 基础用法示例
- 高级功能示例
- 最佳实践

## 参考项目

本 SDK 的设计参考了以下优秀项目：

1. **go-moonshot** - Builder 模式、错误处理
2. **deepseek-go** - 项目结构、接口设计
3. **go-openai** - 流式响应、配置管理

## 未来计划

### 短期

- [ ] 添加更多单元测试
- [ ] 性能基准测试
- [ ] 完善文档

### 中期

- [ ] 向量服务支持
- [ ] 图像生成服务
- [ ] 批处理 API

### 长期

- [ ] gRPC 支持
- [ ] 分布式追踪
- [ ] 指标收集

## 贡献指南

欢迎贡献！在添加新功能时，请遵循以下原则：

1. 保持接口简洁
2. 提供完整的文档
3. 编写测试用例
4. 遵循现有代码风格
5. 考虑向后兼容性

## 总结

360智脑 Go SDK 采用了现代 Go 语言的最佳实践，提供了：

- ✅ 清晰的架构分层
- ✅ 灵活的配置系统
- ✅ 完善的错误处理
- ✅ 良好的扩展性
- ✅ 优秀的开发体验

这使得 SDK 既易于使用，又便于维护和扩展。
