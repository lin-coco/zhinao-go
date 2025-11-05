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
	usage := Usage{
		PromptTokens:     10,
		CompletionTokens: 20,
		TotalTokens:      30,
	}

	if usage.TotalTokens != usage.PromptTokens+usage.CompletionTokens {
		t.Error("Total tokens should equal prompt + completion")
	}
}
