package zhinao

import (
	"os"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		opts    []Option
		wantErr bool
		setup   func()
	}{
		{
			name:    "valid api key",
			apiKey:  "test-api-key",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "empty api key without env",
			apiKey:  "",
			opts:    nil,
			wantErr: true,
			setup: func() {
				// 确保环境变量未设置
				os.Unsetenv(EnvAPIKey)
			},
		},
		{
			name:   "with timeout option",
			apiKey: "test-api-key",
			opts: []Option{
				WithTimeout(30 * time.Second),
			},
			wantErr: false,
		},
		{
			name:   "with base url option",
			apiKey: "test-api-key",
			opts: []Option{
				WithBaseURL("https://custom-api.360.cn/v1"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			client, err := NewClient(tt.apiKey, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestNewClientFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envKey  string
		wantErr bool
	}{
		{
			name:    "with env key set",
			envKey:  "test-api-key-from-env",
			wantErr: false,
		},
		{
			name:    "without env key",
			envKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.envKey != "" {
				os.Setenv(EnvAPIKey, tt.envKey)
				defer os.Unsetenv(EnvAPIKey)
			} else {
				os.Unsetenv(EnvAPIKey)
			}

			client, err := NewClientFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClientFromEnv() returned nil client")
			}
		})
	}
}

func TestClientWithOptions(t *testing.T) {
	apiKey := "test-api-key"

	t.Run("with timeout", func(t *testing.T) {
		timeout := 45 * time.Second
		client, err := NewClient(apiKey, WithTimeout(timeout))
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.config.Timeout != timeout {
			t.Errorf("Timeout = %v, want %v", client.config.Timeout, timeout)
		}
	})

	t.Run("with base url", func(t *testing.T) {
		baseURL := "https://custom.api.com/v1"
		client, err := NewClient(apiKey, WithBaseURL(baseURL))
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.config.BaseURL != baseURL {
			t.Errorf("BaseURL = %v, want %v", client.config.BaseURL, baseURL)
		}
	})

	t.Run("with headers", func(t *testing.T) {
		headers := map[string]string{
			"X-Custom-Header": "test-value",
		}
		client, err := NewClient(apiKey, WithHeaders(headers))
		if err != nil {
			t.Fatalf("NewClient() error = %v", err)
		}
		if client.config.Headers["X-Custom-Header"] != "test-value" {
			t.Errorf("Headers not set correctly")
		}
	})
}

func TestGetConfig(t *testing.T) {
	apiKey := "test-api-key"
	client, err := NewClient(apiKey)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	config := client.GetConfig()
	if config.APIKey != apiKey {
		t.Errorf("GetConfig().APIKey = %v, want %v", config.APIKey, apiKey)
	}
}

func TestClientEnvFallback(t *testing.T) {
	envKey := "test-env-api-key"
	os.Setenv(EnvAPIKey, envKey)
	defer os.Unsetenv(EnvAPIKey)

	// 测试空字符串时自动使用环境变量
	client, err := NewClient("")
	if err != nil {
		t.Fatalf("NewClient(\"\") should use env var, got error: %v", err)
	}
	if client.config.APIKey != envKey {
		t.Errorf("Expected API key from env: %v, got: %v", envKey, client.config.APIKey)
	}
}
