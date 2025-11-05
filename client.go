package zhinao

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// EnvAPIKey 环境变量名称
	EnvAPIKey = "ZHINAO_API_KEY"
)

// Client 360智脑 API 客户端
type Client struct {
	config   *Config
	httpDoer HTTPDoer
}

// NewClient 创建新的客户端实例
//
// 参数:
//   - apiKey: 360智脑 API 密钥（可选，如果为空则从环境变量 ZHINAO_API_KEY 读取）
//   - opts: 可选的配置选项
//
// 返回:
//   - *Client: 客户端实例
//   - error: 如果配置无效则返回错误
//
// 示例:
//
//	// 使用 API Key
//	client, err := zhinao.NewClient("your-api-key")
//
//	// 使用环境变量（设置 ZHINAO_API_KEY）
//	client, err := zhinao.NewClient("")
//
//	// 使用自定义配置
//	client, err := zhinao.NewClient(
//	    "your-api-key",
//	    zhinao.WithTimeout(30*time.Second),
//	)
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	// 如果 apiKey 为空，尝试从环境变量读取
	if apiKey == "" {
		apiKey = os.Getenv(EnvAPIKey)
	}

	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	// 创建默认配置
	config := &Config{
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL,
		Timeout: DefaultTimeout,
	}

	// 应用配置选项
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 创建 HTTP Doer
	var httpDoer HTTPDoer
	if config.HTTPDoer != nil {
		// 用户提供了自定义的 HTTPDoer，直接使用
		httpDoer = config.HTTPDoer
	} else {
		// 使用配置创建标准的 http.Client
		httpDoer = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		config:   config,
		httpDoer: httpDoer,
	}, nil
}

// NewClientFromEnv 从环境变量创建客户端
//
// 这是一个便捷方法，等同于 NewClient("")
//
// 示例:
//
//	// 确保设置了环境变量 ZHINAO_API_KEY
//	client, err := zhinao.NewClientFromEnv()
func NewClientFromEnv(opts ...Option) (*Client, error) {
	return NewClient("", opts...)
}

// GetConfig 获取客户端配置（只读）
func (c *Client) GetConfig() Config {
	return *c.config
}

// doRequest 执行 HTTP 请求并返回响应
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// 执行请求
	resp, err := c.httpDoer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API error (status %d): failed to read error response", resp.StatusCode)
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	return resp, nil
}

// buildRequest 构建 HTTP 请求
func (c *Client) buildRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	// 设置自定义请求头
	for k, v := range c.config.Headers {
		req.Header.Set(k, v)
	}

	return req, nil
}
