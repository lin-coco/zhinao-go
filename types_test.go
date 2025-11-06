package zhinao

import (
	"testing"
)

func TestChatRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatCompletionRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty model",
			req: &ChatCompletionRequest{
				Model: "",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			wantErr: true,
		},
		{
			name: "empty messages",
			req: &ChatCompletionRequest{
				Model:    "360gpt-turbo",
				Messages: []ChatCompletionMessage{},
			},
			wantErr: true,
		},
		{
			name: "nil messages",
			req: &ChatCompletionRequest{
				Model:    "360gpt-turbo",
				Messages: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid temperature - too low",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: -0.1,
			},
			wantErr: true,
		},
		{
			name: "invalid temperature - too high",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: 2.1,
			},
			wantErr: true,
		},
		{
			name: "valid temperature",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: 0.7,
			},
			wantErr: false,
		},
		{
			name: "invalid top_p - too low",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
				TopP: -0.1,
			},
			wantErr: true,
		},
		{
			name: "invalid top_p - too high",
			req: &ChatCompletionRequest{
				Model: "360gpt-turbo",
				Messages: []ChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
				TopP: 1.1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageValidation(t *testing.T) {
	t.Run("valid message", func(t *testing.T) {
		msg := ChatCompletionMessage{
			Role:    "user",
			Content: "Hello",
		}
		if msg.Role == "" || msg.Content == "" {
			t.Error("ChatCompletionMessage should be valid")
		}
	})

	t.Run("message with tool calls", func(t *testing.T) {
		msg := ChatCompletionMessage{
			Role: "assistant",
			ToolCalls: []ToolCall{
				{
					ID:   "call_123",
					Type: "function",
					Function: ToolCallFunction{
						Name:      "get_weather",
						Arguments: `{"city": "Beijing"}`,
					},
				},
			},
		}
		if len(msg.ToolCalls) != 1 {
			t.Error("Tool calls should be set")
		}
	})
}

func TestUsageInfo(t *testing.T) {
	t.Run("basic usage info", func(t *testing.T) {
		usage := Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		}

		if usage.TotalTokens != usage.PromptTokens+usage.CompletionTokens {
			t.Error("Total tokens should equal prompt + completion")
		}
	})

	t.Run("usage with cache info", func(t *testing.T) {
		usage := Usage{
			PromptTokens:          425,
			CompletionTokens:      67,
			TotalTokens:           492,
			PromptCacheHitTokens:  320,
			PromptCacheMissTokens: 105,
		}

		// 验证总token数等于提示词token数加上完成token数
		if usage.TotalTokens != usage.PromptTokens+usage.CompletionTokens {
			t.Errorf("Total tokens (%d) should equal prompt tokens (%d) + completion tokens (%d)",
				usage.TotalTokens, usage.PromptTokens, usage.CompletionTokens)
		}

		// 验证缓存命中token数和未命中token数之和等于提示词token数
		if usage.PromptCacheHitTokens+usage.PromptCacheMissTokens != usage.PromptTokens {
			t.Errorf("Cache hit tokens (%d) + cache miss tokens (%d) should equal prompt tokens (%d)",
				usage.PromptCacheHitTokens, usage.PromptCacheMissTokens, usage.PromptTokens)
		}
	})

	t.Run("usage without cache (backward compatibility)", func(t *testing.T) {
		usage := Usage{
			PromptTokens:          20,
			CompletionTokens:      14,
			TotalTokens:           34,
			PromptCacheHitTokens:  0,
			PromptCacheMissTokens: 20,
		}

		// 当没有缓存命中时，所有提示词token都应该是未命中的
		if usage.PromptCacheMissTokens != usage.PromptTokens {
			t.Errorf("When no cache hits, cache miss tokens (%d) should equal prompt tokens (%d)",
				usage.PromptCacheMissTokens, usage.PromptTokens)
		}

		if usage.PromptCacheHitTokens != 0 {
			t.Errorf("Expected no cache hits (0), got %d", usage.PromptCacheHitTokens)
		}
	})
}

// TestFunctionParameters 测试 FunctionParameters 结构体
func TestFunctionParameters(t *testing.T) {
	t.Run("basic function parameters", func(t *testing.T) {
		params := FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "城市名称",
				},
			},
			Required: []string{"location"},
		}

		if params.Type != "object" {
			t.Errorf("Expected type 'object', got '%s'", params.Type)
		}

		if len(params.Required) != 1 || params.Required[0] != "location" {
			t.Errorf("Expected required ['location'], got %v", params.Required)
		}
	})

	t.Run("function parameters from API example", func(t *testing.T) {
		// 基于官方文档示例
		params := FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"from": map[string]interface{}{
					"type":        "string",
					"description": "出发地",
				},
				"to": map[string]interface{}{
					"type":        "string",
					"description": "目的地",
				},
			},
			Required: []string{"from", "to"},
		}

		if params.Type != "object" {
			t.Errorf("Expected type 'object', got '%s'", params.Type)
		}

		if len(params.Properties) != 2 {
			t.Errorf("Expected 2 properties, got %d", len(params.Properties))
		}

		if len(params.Required) != 2 {
			t.Errorf("Expected 2 required fields, got %d", len(params.Required))
		}
	})
}

// TestToolChoiceFunction 测试 ToolChoiceFunction 结构体
func TestToolChoiceFunction(t *testing.T) {
	t.Run("basic tool choice function", func(t *testing.T) {
		tcf := ToolChoiceFunction{
			Name: "getTrain",
		}

		if tcf.Name != "getTrain" {
			t.Errorf("Expected name 'getTrain', got '%s'", tcf.Name)
		}
	})
}

// TestToolChoice 测试 ToolChoice 结构体
func TestToolChoice(t *testing.T) {
	t.Run("tool choice with function", func(t *testing.T) {
		toolChoice := ToolChoice{
			Type: "function",
			Function: ToolChoiceFunction{
				Name: "my_function",
			},
		}

		if toolChoice.Type != "function" {
			t.Errorf("Expected type 'function', got '%s'", toolChoice.Type)
		}

		if toolChoice.Function.Name != "my_function" {
			t.Errorf("Expected function name 'my_function', got '%s'", toolChoice.Function.Name)
		}
	})
}

// TestToolFunctionWithFunctionParameters 测试 ToolFunction 使用 FunctionParameters
func TestToolFunctionWithFunctionParameters(t *testing.T) {
	t.Run("tool function with FunctionParameters", func(t *testing.T) {
		params := &FunctionParameters{
			Type: "object",
			Properties: map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "查询景点的地名",
				},
			},
			Required: []string{"location"},
		}

		toolFunc := ToolFunction{
			Name:        "getSights",
			Description: "查询景点",
			Parameters:  params,
		}

		if toolFunc.Name != "getSights" {
			t.Errorf("Expected name 'getSights', got '%s'", toolFunc.Name)
		}

		if toolFunc.Parameters == nil {
			t.Error("Parameters should not be nil")
		}

		if toolFunc.Parameters.Type != "object" {
			t.Errorf("Expected parameters type 'object', got '%s'", toolFunc.Parameters.Type)
		}
	})
}

// TestCompleteToolDefinition 测试完整的工具定义（基于官方文档）
func TestCompleteToolDefinition(t *testing.T) {
	t.Run("complete tool from API example", func(t *testing.T) {
		tool := Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        "getTrain",
				Description: "查询火车票",
				Parameters: &FunctionParameters{
					Type: "object",
					Properties: map[string]interface{}{
						"from": map[string]interface{}{
							"type":        "string",
							"description": "出发地",
						},
						"to": map[string]interface{}{
							"type":        "string",
							"description": "目的地",
						},
					},
					Required: []string{"from", "to"},
				},
			},
		}

		if tool.Type != "function" {
			t.Errorf("Expected type 'function', got '%s'", tool.Type)
		}

		if tool.Function.Name != "getTrain" {
			t.Errorf("Expected function name 'getTrain', got '%s'", tool.Function.Name)
		}

		if tool.Function.Description != "查询火车票" {
			t.Errorf("Expected description '查询火车票', got '%s'", tool.Function.Description)
		}

		if tool.Function.Parameters == nil {
			t.Fatal("Parameters should not be nil")
		}

		if len(tool.Function.Parameters.Required) != 2 {
			t.Errorf("Expected 2 required parameters, got %d", len(tool.Function.Parameters.Required))
		}
	})
}

// TestChatCompletionStreamResponse 测试流式响应结构体
func TestChatCompletionStreamResponse(t *testing.T) {
	t.Run("stream response with usage", func(t *testing.T) {
		response := ChatCompletionStreamResponse{
			ID:      "chatcmpl-123",
			Object:  "chat.completion.chunk",
			Created: 1234567890,
			Model:   "360gpt-turbo",
			Choices: []ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: Delta{
						Role:    "assistant",
						Content: "Hello",
					},
					FinishReason: "stop",
				},
			},
			Usage: &Usage{
				PromptTokens:          425,
				CompletionTokens:      67,
				TotalTokens:           492,
				PromptCacheHitTokens:  320,
				PromptCacheMissTokens: 105,
			},
		}

		if response.ID != "chatcmpl-123" {
			t.Errorf("Expected ID 'chatcmpl-123', got '%s'", response.ID)
		}

		if response.Model != "360gpt-turbo" {
			t.Errorf("Expected model '360gpt-turbo', got '%s'", response.Model)
		}

		if len(response.Choices) != 1 {
			t.Errorf("Expected 1 choice, got %d", len(response.Choices))
		}

		// 验证 Usage 字段不为 nil
		if response.Usage == nil {
			t.Fatal("Expected Usage to be non-nil")
		}

		// 验证 Usage 字段的值
		if response.Usage.TotalTokens != 492 {
			t.Errorf("Expected total tokens 492, got %d", response.Usage.TotalTokens)
		}

		if response.Usage.PromptCacheHitTokens != 320 {
			t.Errorf("Expected cache hit tokens 320, got %d", response.Usage.PromptCacheHitTokens)
		}
	})

	t.Run("stream response without usage", func(t *testing.T) {
		response := ChatCompletionStreamResponse{
			ID:      "chatcmpl-456",
			Object:  "chat.completion.chunk",
			Created: 1234567890,
			Model:   "360gpt-turbo",
			Choices: []ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: Delta{
						Content: "World",
					},
				},
			},
			// Usage 字段为 nil，表示此块不包含 usage 信息
			Usage: nil,
		}

		if response.ID != "chatcmpl-456" {
			t.Errorf("Expected ID 'chatcmpl-456', got '%s'", response.ID)
		}

		// 验证 Usage 为 nil
		if response.Usage != nil {
			t.Errorf("Expected Usage to be nil, got %v", response.Usage)
		}
	})

	t.Run("stream response with tool calls", func(t *testing.T) {
		response := ChatCompletionStreamResponse{
			ID:      "chatcmpl-789",
			Object:  "chat.completion.chunk",
			Created: 1234567890,
			Model:   "360gpt-turbo",
			Choices: []ChatCompletionStreamChoice{
				{
					Index: 0,
					Delta: Delta{
						Role: "assistant",
						ToolCalls: []ToolCall{
							{
								ID:   "call_123",
								Type: "function",
								Function: ToolCallFunction{
									Name:      "get_weather",
									Arguments: `{"city": "Beijing"}`,
								},
							},
						},
					},
					FinishReason: "tool_calls",
				},
			},
			Usage: &Usage{
				PromptTokens:     50,
				CompletionTokens: 30,
				TotalTokens:      80,
			},
		}

		if len(response.Choices) != 1 {
			t.Fatalf("Expected 1 choice, got %d", len(response.Choices))
		}

		if len(response.Choices[0].Delta.ToolCalls) != 1 {
			t.Errorf("Expected 1 tool call, got %d", len(response.Choices[0].Delta.ToolCalls))
		}

		if response.Choices[0].FinishReason != "tool_calls" {
			t.Errorf("Expected finish reason 'tool_calls', got '%s'", response.Choices[0].FinishReason)
		}

		// 验证 Usage 不为 nil
		if response.Usage == nil {
			t.Fatal("Expected Usage to be non-nil")
		}

		if response.Usage.TotalTokens != 80 {
			t.Errorf("Expected total tokens 80, got %d", response.Usage.TotalTokens)
		}
	})
}
