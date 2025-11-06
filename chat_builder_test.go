package zhinao

import (
	"testing"
)

func TestChatBuilder(t *testing.T) {
	t.Run("basic build", func(t *testing.T) {
		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddUserMessage("Hello").
			Build()

		if req.Model != "360gpt-turbo" {
			t.Errorf("Model = %v, want %v", req.Model, "360gpt-turbo")
		}
		if len(req.Messages) != 1 {
			t.Errorf("Messages length = %v, want %v", len(req.Messages), 1)
		}
		if req.Messages[0].Role != "user" {
			t.Errorf("Message role = %v, want %v", req.Messages[0].Role, "user")
		}
		if req.Messages[0].Content != "Hello" {
			t.Errorf("Message content = %v, want %v", req.Messages[0].Content, "Hello")
		}
	})

	t.Run("multiple messages", func(t *testing.T) {
		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddSystemMessage("You are a helper").
			AddUserMessage("Hello").
			AddAssistantMessage("Hi there!").
			AddUserMessage("How are you?").
			Build()

		if len(req.Messages) != 4 {
			t.Errorf("Messages length = %v, want %v", len(req.Messages), 4)
		}

		expectedRoles := []string{"system", "user", "assistant", "user"}
		for i, msg := range req.Messages {
			if msg.Role != expectedRoles[i] {
				t.Errorf("Message[%d] role = %v, want %v", i, msg.Role, expectedRoles[i])
			}
		}
	})

	t.Run("with parameters", func(t *testing.T) {
		temperature := float32(0.7)
		maxTokens := 1000
		topP := float32(0.9)

		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddUserMessage("Test").
			SetTemperature(temperature).
			SetMaxTokens(maxTokens).
			SetTopP(topP).
			Build()

		if req.Temperature != temperature {
			t.Errorf("Temperature = %v, want %v", req.Temperature, temperature)
		}
		if req.MaxTokens != maxTokens {
			t.Errorf("MaxTokens = %v, want %v", req.MaxTokens, maxTokens)
		}
		if req.TopP != topP {
			t.Errorf("TopP = %v, want %v", req.TopP, topP)
		}
	})

	t.Run("add messages", func(t *testing.T) {
		messages := []ChatCompletionMessage{
			{Role: "system", Content: "System message"},
			{Role: "user", Content: "User message"},
		}

		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddMessages(messages).
			Build()

		if len(req.Messages) != 2 {
			t.Errorf("Messages length = %v, want %v", len(req.Messages), 2)
		}
	})

	t.Run("with tools", func(t *testing.T) {
		tool := Tool{
			Type: "function",
			Function: ToolFunction{
				Name:        "get_weather",
				Description: "Get weather info",
			},
		}

		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddUserMessage("Test").
			AddTool(tool).
			Build()

		if len(req.Tools) != 1 {
			t.Errorf("Tools length = %v, want %v", len(req.Tools), 1)
		}
		if req.Tools[0].Function.Name != "get_weather" {
			t.Errorf("Tool name = %v, want %v", req.Tools[0].Function.Name, "get_weather")
		}
	})

	t.Run("with multiple tools", func(t *testing.T) {
		tools := []Tool{
			{
				Type: "function",
				Function: ToolFunction{
					Name: "get_weather",
				},
			},
			{
				Type: "function",
				Function: ToolFunction{
					Name: "get_time",
				},
			},
		}

		req := NewChatCompletionBuilder().
			SetModel("360gpt-turbo").
			AddUserMessage("Test").
			AddTools(tools).
			Build()

		if len(req.Tools) != 2 {
			t.Errorf("Tools length = %v, want %v", len(req.Tools), 2)
		}
	})
}

func TestChatBuilderChaining(t *testing.T) {
	// 测试每次 Build() 返回的是同一个请求对象
	// Builder 模式修改的是内部请求对象
	builder := NewChatCompletionBuilder()

	builder.SetModel("model-1").AddUserMessage("Message 1")
	req1 := builder.Build()

	builder.SetModel("model-2").AddUserMessage("Message 2")
	req2 := builder.Build()

	// 由于 Builder 修改的是同一个内部对象，req1 和 req2 都指向同一个请求
	if req1 != req2 {
		t.Error("Builder should return the same request object")
	}

	// 验证最后的状态
	if req2.Model != "model-2" {
		t.Errorf("Model = %v, want model-2", req2.Model)
	}
	if len(req2.Messages) != 2 {
		t.Errorf("Messages length = %v, want 2", len(req2.Messages))
	}
}

func TestChatBuilderIndependence(t *testing.T) {
	// 测试不同 Builder 实例的独立性
	builder1 := NewChatCompletionBuilder().
		SetModel("model-1").
		AddUserMessage("Message 1")

	builder2 := NewChatCompletionBuilder().
		SetModel("model-2").
		AddUserMessage("Message 2")

	req1 := builder1.Build()
	req2 := builder2.Build()

	// 不同 Builder 应该创建不同的请求
	if req1 == req2 {
		t.Error("Different builders should create different requests")
	}
	if req1.Model == req2.Model {
		t.Error("Different builders should have different models")
	}
}
