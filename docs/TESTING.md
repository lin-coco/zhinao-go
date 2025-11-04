# 测试指南

本文档介绍 360智脑 Go SDK 的测试策略和最佳实践。

## 测试概览

项目采用全面的单元测试策略，确保代码质量和稳定性。所有核心功能都有对应的测试用例。

### 测试覆盖范围

- ✅ **客户端管理** - 客户端创建、配置、选项处理
- ✅ **环境变量支持** - API 密钥环境变量读取
- ✅ **Builder 模式** - 请求构建器的链式调用和独立性
- ✅ **错误处理** - API 错误、限流错误、验证错误
- ✅ **请求验证** - 参数验证、边界条件检查
- ✅ **类型定义** - 消息、工具调用等类型验证

## 运行测试

### 基本命令

```bash
# 运行所有测试
go test -v ./...

# 仅运行单元测试（不包括 examples）
go test -v .

# 运行特定测试
go test -v -run TestChatBuilder

# 运行特定包的测试
go test -v ./internal/http
```

### 使用 Makefile

项目提供了便捷的 Makefile 命令：

```bash
# 运行所有测试
make test

# 仅运行单元测试
make test-unit

# 生成覆盖率报告
make test-coverage

# 运行代码检查
make lint

# 格式化代码
make fmt
```

## 测试结构

### 测试文件组织

```
zhinao-go/
├── client_test.go        # 客户端测试
├── chat_builder_test.go  # Builder 模式测试
├── errors_test.go        # 错误处理测试
├── types_test.go         # 类型定义测试
└── internal/
    └── http/
        └── *_test.go     # HTTP 客户端测试
```

### 测试命名规范

测试函数遵循以下命名规范：

```go
func TestFunctionName(t *testing.T)           // 测试单个函数
func TestFunctionName_Scenario(t *testing.T)  // 测试特定场景
```

使用子测试组织相关测试用例：

```go
func TestClient(t *testing.T) {
    t.Run("valid api key", func(t *testing.T) {
        // 测试代码
    })
    
    t.Run("empty api key", func(t *testing.T) {
        // 测试代码
    })
}
```

## 测试示例

### 客户端测试

```go
func TestNewClient(t *testing.T) {
    client, err := NewClient("test-api-key")
    if err != nil {
        t.Fatalf("NewClient() error = %v", err)
    }
    if client == nil {
        t.Error("NewClient() returned nil client")
    }
}
```

### Builder 测试

```go
func TestChatBuilder(t *testing.T) {
    req := NewChatBuilder().
        SetModel("360gpt-turbo").
        AddUserMessage("Hello").
        Build()
    
    if req.Model != "360gpt-turbo" {
        t.Errorf("Model = %v, want %v", req.Model, "360gpt-turbo")
    }
}
```

### 错误处理测试

```go
func TestAPIError(t *testing.T) {
    err := &APIError{
        StatusCode: 400,
        Message:    "Bad request",
    }
    
    if !strings.Contains(err.Error(), "400") {
        t.Error("Error message should contain status code")
    }
}
```

## 覆盖率报告

### 生成覆盖率报告

```bash
# 使用 Makefile
make test-coverage

# 或手动执行
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

生成的 `coverage.html` 文件可以在浏览器中打开查看详细的覆盖率信息。

### 查看覆盖率统计

```bash
# 查看总体覆盖率
go test -cover ./...

# 查看详细覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## 测试最佳实践

### 1. 测试独立性

每个测试应该是独立的，不依赖其他测试的执行顺序：

```go
func TestIndependent(t *testing.T) {
    // 每个测试创建自己的数据
    client, _ := NewClient("test-key")
    // 测试逻辑
}
```

### 2. 使用表格驱动测试

对于多个类似的测试场景，使用表格驱动测试：

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "valid", false},
        {"empty input", "", true},
        {"invalid input", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 3. 清理资源

使用 `defer` 确保测试资源被正确清理：

```go
func TestWithCleanup(t *testing.T) {
    client, _ := NewClient("test-key")
    defer client.Close() // 如果有 Close 方法
    
    // 测试逻辑
}
```

### 4. 测试错误情况

不仅测试正常情况，也要测试错误情况：

```go
func TestErrorHandling(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // 测试成功情况
    })
    
    t.Run("error case", func(t *testing.T) {
        // 测试错误情况
    })
}
```

### 5. 使用辅助函数

提取通用的测试设置代码到辅助函数：

```go
func newTestClient(t *testing.T) *Client {
    client, err := NewClient("test-api-key")
    if err != nil {
        t.Fatalf("Failed to create test client: %v", err)
    }
    return client
}

func TestWithHelper(t *testing.T) {
    client := newTestClient(t)
    // 使用 client 进行测试
}
```

## 持续集成

### GitHub Actions 配置示例

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        run: make test
      
      - name: Generate coverage
        run: make test-coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
```

## 贡献测试

如果你想为项目贡献测试用例：

1. **确保测试通过** - 所有新测试都应该通过
2. **保持覆盖率** - 新代码应该有对应的测试
3. **遵循规范** - 遵循项目的测试命名和组织规范
4. **添加文档** - 为复杂的测试场景添加注释

### 提交测试的检查清单

- [ ] 测试命名清晰且符合规范
- [ ] 使用子测试组织相关用例
- [ ] 测试覆盖了正常和异常情况
- [ ] 所有测试都能独立运行
- [ ] 添加了必要的注释
- [ ] 运行 `make test` 确保所有测试通过
- [ ] 运行 `make lint` 确保代码风格正确

## 故障排查

### 常见问题

**问题**: 测试失败，错误信息为 "panic: runtime error"

**解决**: 检查是否正确初始化了所有必需的字段，特别是指针类型。

**问题**: 环境变量相关测试失败

**解决**: 确保在测试中正确设置和清理环境变量：
```go
os.Setenv("KEY", "value")
defer os.Unsetenv("KEY")
```

**问题**: 测试超时

**解决**: 检查是否有无限循环或阻塞操作，考虑添加超时：
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## 参考资源

- [Go 测试官方文档](https://golang.org/pkg/testing/)
- [表格驱动测试](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go 测试最佳实践](https://go.dev/doc/tutorial/add-a-test)
