package zhinao

import (
	"fmt"
	"net/url"
	"time"

	"github.com/lin-coco/zhinao-go/internal/http"
)

const (
	// DefaultBaseURL 默认 API 基础 URL
	DefaultBaseURL = "https://api.360.cn/v1"

	// DefaultTimeout 默认超时时间
	DefaultTimeout = 60 * time.Second

	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 3

	// DefaultRetryDelay 默认重试延迟
	DefaultRetryDelay = 1 * time.Second
)

// Config 客户端配置
type Config struct {
	APIKey     string            // API 密钥（必需）
	BaseURL    string            // API 基础 URL
	Timeout    time.Duration     // 请求超时时间
	MaxRetries int               // 最大重试次数
	RetryDelay time.Duration     // 重试延迟
	HTTPClient http.Client       // 自定义 HTTP 客户端
	UserAgent  string            // 自定义 User-Agent
	Headers    map[string]string // 自定义请求头
}

// Option 配置选项函数
type Option func(*Config) error

// WithBaseURL 设置 API 基础 URL
//
// 示例:
//
//	client, err := NewClient(apiKey, WithBaseURL("https://custom-api.360.cn"))
func WithBaseURL(baseURL string) Option {
	return func(c *Config) error {
		if _, err := url.Parse(baseURL); err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
		c.BaseURL = baseURL
		return nil
	}
}

// WithTimeout 设置请求超时时间
//
// 示例:
//
//	client, err := NewClient(apiKey, WithTimeout(30*time.Second))
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) error {
		if timeout < 0 {
			return fmt.Errorf("timeout must be positive")
		}
		c.Timeout = timeout
		return nil
	}
}

// WithRetry 设置重试配置
//
// 参数:
//   - maxRetries: 最大重试次数
//   - retryDelay: 基础重试延迟（使用指数退避策略）
//
// 示例:
//
//	client, err := NewClient(apiKey, WithRetry(5, 2*time.Second))
func WithRetry(maxRetries int, retryDelay time.Duration) Option {
	return func(c *Config) error {
		if maxRetries < 0 {
			return fmt.Errorf("max retries must be non-negative")
		}
		if retryDelay < 0 {
			return fmt.Errorf("retry delay must be non-negative")
		}
		c.MaxRetries = maxRetries
		c.RetryDelay = retryDelay
		return nil
	}
}

// WithHTTPClient 设置自定义 HTTP 客户端
//
// 这允许你使用自定义的 HTTP 客户端实现，例如用于特殊的代理配置或其他 HTTP 库
//
// 示例:
//
//	customClient := &MyCustomHTTPClient{}
//	client, err := NewClient(apiKey, WithHTTPClient(customClient))
func WithHTTPClient(client http.Client) Option {
	return func(c *Config) error {
		if client == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}
		c.HTTPClient = client
		return nil
	}
}

// WithUserAgent 设置自定义 User-Agent
//
// 示例:
//
//	client, err := NewClient(apiKey, WithUserAgent("MyApp/1.0"))
func WithUserAgent(userAgent string) Option {
	return func(c *Config) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithHeaders 设置自定义请求头
//
// 示例:
//
//	headers := map[string]string{
//	    "X-Custom-Header": "value",
//	}
//	client, err := NewClient(apiKey, WithHeaders(headers))
func WithHeaders(headers map[string]string) Option {
	return func(c *Config) error {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
		return nil
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base URL cannot be empty")
	}
	if c.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}
	return nil
}
