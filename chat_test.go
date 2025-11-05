package zhinao

import (
	"context"
	"testing"

	"github.com/lin-coco/zhinao-go/internal/test"
)

// setupZhinaoTestServer 设置测试服务器和客户端
func setupZhinaoTestServer() (client *Client, server *test.ServerTest, teardown func()) {
	server = test.NewTestServer()
	ts := server.ZhinaoTestServer()
	ts.Start()
	teardown = ts.Close

	// 创建指向测试服务器的客户端
	client, _ = NewClient(test.GetTestToken(), WithBaseURL(ts.URL))
	return
}

// TestChatCompletion_Mock 测试聊天补全（使用 Mock Server）
func TestChatCompletion_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册模拟的聊天补全响应
	server.RegisterHandler("/chat/completions", test.ChatCompletionResponse())

	// 执行聊天补全请求
	req := &ChatRequest{
		Model: Model360GPTTurbo,
		Messages: []Message{
			{Role: RoleUser, Content: "Hello!"},
		},
	}

	resp, err := client.CreateCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCompletion failed: %v", err)
	}

	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}

	if resp.ID != "chatcmpl-test-123" {
		t.Errorf("Expected ID 'chatcmpl-test-123', got '%s'", resp.ID)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}

	if resp.Choices[0].Message.Role != RoleAssistant {
		t.Errorf("Expected role 'assistant', got '%s'", resp.Choices[0].Message.Role)
	}

	if resp.Choices[0].Message.Content == "" {
		t.Error("Response content is empty")
	}
}

// TestChatCompletion_WithBuilder_Mock 测试使用 Builder 的聊天补全
func TestChatCompletion_WithBuilder_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	server.RegisterHandler("/chat/completions", test.ChatCompletionResponse())

	// 使用 Builder 构建请求
	req := NewChatBuilder().
		SetModel(Model360GPTTurbo).
		AddUserMessage("测试消息").
		AddSystemMessage("你是一个helpful的助手").
		SetTemperature(0.7).
		SetMaxTokens(100).
		Build()

	resp, err := client.CreateCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCompletion failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}
}

// TestChatCompletion_InvalidAuth_Mock 测试无效的认证
func TestChatCompletion_InvalidAuth_Mock(t *testing.T) {
	server := test.NewTestServer()
	ts := server.ZhinaoTestServer()
	ts.Start()
	defer ts.Close()

	// 使用错误的 API Key
	client, _ := NewClient("invalid-api-key", WithBaseURL(ts.URL))

	server.RegisterHandler("/chat/completions", test.ChatCompletionResponse())

	req := &ChatRequest{
		Model: Model360GPTTurbo,
		Messages: []Message{
			{Role: RoleUser, Content: "Hello!"},
		},
	}

	_, err := client.CreateCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error for invalid API key, got nil")
	}

	// 验证返回的是 API 错误
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != 401 {
		t.Errorf("Expected status code 401, got %d", apiErr.StatusCode)
	}
}

// TestChatCompletion_RateLimit_Mock 测试限流错误
func TestChatCompletion_RateLimit_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册限流响应
	server.RegisterHandler("/chat/completions", test.RateLimitResponse(60))

	req := &ChatRequest{
		Model: Model360GPTTurbo,
		Messages: []Message{
			{Role: RoleUser, Content: "Hello!"},
		},
	}

	_, err := client.CreateCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("Expected rate limit error, got nil")
	}

	// 验证是 API 错误
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != 429 {
		t.Errorf("Expected status code 429, got %d", apiErr.StatusCode)
	}
}

// TestChatCompletion_ServerError_Mock 测试服务器错误
func TestChatCompletion_ServerError_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册服务器错误响应
	server.RegisterHandler("/chat/completions",
		test.ErrorResponse(500, "Internal server error", "server_error", "internal_error"))

	req := &ChatRequest{
		Model: Model360GPTTurbo,
		Messages: []Message{
			{Role: RoleUser, Content: "Hello!"},
		},
	}

	_, err := client.CreateCompletion(context.Background(), req)
	if err == nil {
		t.Fatal("Expected server error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", apiErr.StatusCode)
	}
}

// TestChatCompletion_MultipleMessages_Mock 测试多轮对话
func TestChatCompletion_MultipleMessages_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	server.RegisterHandler("/chat/completions", test.ChatCompletionResponse())

	// 构建多轮对话
	req := NewChatBuilder().
		SetModel(Model360GPTTurbo).
		AddSystemMessage("你是一个helpful的助手").
		AddUserMessage("你好").
		AddAssistantMessage("你好！有什么我可以帮助你的吗？").
		AddUserMessage("今天天气怎么样？").
		Build()

	resp, err := client.CreateCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCompletion failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	// 验证请求包含了所有消息
	if len(req.Messages) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(req.Messages))
	}
}

// TestChatCompletion_WithTools_Mock 测试带工具调用的请求
func TestChatCompletion_WithTools_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	server.RegisterHandler("/chat/completions", test.ChatCompletionResponse())

	// 定义工具
	tool := Tool{
		Type: "function",
		Function: ToolFunction{
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
	}

	req := NewChatBuilder().
		SetModel(Model360GPTTurbo).
		AddUserMessage("北京今天天气怎么样？").
		AddTool(tool).
		Build()

	resp, err := client.CreateCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateCompletion failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Response is nil")
	}

	// 验证工具被正确添加
	if len(req.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(req.Tools))
	}
}
