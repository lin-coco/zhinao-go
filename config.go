package zhinao

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// HTTPDoer 是执行 HTTP 请求的接口
// 标准库的 *http.Client 实现了此接口
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	// DefaultBaseURL 默认 API 基础 URL
	DefaultBaseURL = "https://api.360.cn/v1"

	// DefaultTimeout 默认超时时间
	DefaultTimeout = 60 * time.Second
)

// Config 客户端配置
type Config struct {
	APIKey   string            // API 密钥（必需）
	BaseURL  string            // API 基础 URL
	Timeout  time.Duration     // 请求超时时间
	HTTPDoer HTTPDoer          // 自定义 HTTP 客户端（实现 Do 方法）
	Headers  map[string]string // 自定义请求头
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

// WithHTTPClient 设置自定义 HTTP 客户端
//
// 接受任何实现了 HTTPDoer 接口的客户端（包括标准库的 *http.Client）
// 当使用此选项时，WithTimeout 配置将被忽略（因为超时应在自定义客户端中设置）
// 但 WithBaseURL 和 WithHeaders 仍然有效
//
// 示例:
//
//	// 使用标准库的 http.Client
//	httpClient := &http.Client{
//	    Timeout: 30 * time.Second,
//	    Transport: &http.Transport{
//	        MaxIdleConns: 100,
//	    },
//	}
//	client, err := NewClient(apiKey, WithHTTPClient(httpClient))
//
//	// 或使用自定义实现
//	customClient := &MyCustomHTTPClient{}
//	client, err := NewClient(apiKey, WithHTTPClient(customClient))
func WithHTTPClient(doer HTTPDoer) Option {
	return func(c *Config) error {
		if doer == nil {
			return fmt.Errorf("HTTP client cannot be nil")
		}
		c.HTTPDoer = doer
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
	return nil
}
