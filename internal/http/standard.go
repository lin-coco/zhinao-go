package http

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// StandardClient 标准 HTTP 客户端实现
type StandardClient struct {
	baseURL     string
	httpClient  *http.Client
	retryConfig *RetryConfig
}

// NewStandardClient 创建标准 HTTP 客户端
func NewStandardClient(baseURL string, timeout time.Duration, maxRetries int, retryDelay time.Duration) Client {
	return &StandardClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
		retryConfig: &RetryConfig{
			MaxRetries: maxRetries,
			RetryDelay: retryDelay,
			Backoff:    ExponentialBackoff,
		},
	}
}

// Post 发送 POST 请求
func (c *StandardClient) Post(ctx context.Context, path string, body, result interface{}, apiKey string) error {
	// 序列化请求体
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 使用重试机制执行请求
	resp, err := executeWithRetry(ctx, c.retryConfig, func() (*http.Response, error) {
		req, err := c.buildRequest(ctx, http.MethodPost, path, bytes.NewReader(jsonData), apiKey)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	return c.parseResponse(resp, result)
}

// Get 发送 GET 请求
func (c *StandardClient) Get(ctx context.Context, path string, result interface{}, apiKey string) error {
	// 使用重试机制执行请求
	resp, err := executeWithRetry(ctx, c.retryConfig, func() (*http.Response, error) {
		req, err := c.buildRequest(ctx, http.MethodGet, path, nil, apiKey)
		if err != nil {
			return nil, err
		}
		return c.httpClient.Do(req)
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	return c.parseResponse(resp, result)
}

// PostStream 发送流式 POST 请求
func (c *StandardClient) PostStream(ctx context.Context, path string, body interface{}, apiKey string) (StreamReader, error) {
	// 序列化请求体
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建请求
	req, err := c.buildRequest(ctx, http.MethodPost, path, bytes.NewReader(jsonData), apiKey)
	if err != nil {
		return nil, err
	}

	// 发送请求（流式请求不使用重试）
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, c.parseError(resp)
	}

	// 返回流读取器
	return &standardStreamReader{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}, nil
}

// buildRequest 构建 HTTP 请求
func (c *StandardClient) buildRequest(ctx context.Context, method, path string, body io.Reader, apiKey string) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", "360AI-Go-SDK/2.0.0")

	return req, nil
}

// parseResponse 解析响应
func (c *StandardClient) parseResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return c.parseErrorFromBody(resp.StatusCode, body)
	}

	// 解析结果
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// parseError 从响应中解析错误
func (c *StandardClient) parseError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("API error (status %d): failed to read error response", resp.StatusCode)
	}
	return c.parseErrorFromBody(resp.StatusCode, body)
}

// parseErrorFromBody 从响应体解析错误
func (c *StandardClient) parseErrorFromBody(statusCode int, body []byte) error {
	// 尝试解析为标准错误格式
	var errResp struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		// 检查是否是限流错误
		if statusCode == http.StatusTooManyRequests {
			// 可以从响应头获取 Retry-After
			return &RateLimitError{
				APIError: &APIError{
					StatusCode: statusCode,
					Type:       errResp.Error.Type,
					Code:       errResp.Error.Code,
					Message:    errResp.Error.Message,
				},
			}
		}

		return &APIError{
			StatusCode: statusCode,
			Type:       errResp.Error.Type,
			Code:       errResp.Error.Code,
			Message:    errResp.Error.Message,
		}
	}

	// 如果无法解析，返回原始错误
	return &APIError{
		StatusCode: statusCode,
		Message:    string(body),
	}
}

// standardStreamReader 标准流读取器实现
type standardStreamReader struct {
	resp   *http.Response
	reader *bufio.Reader
}

// Read 读取下一行数据
func (s *standardStreamReader) Read() ([]byte, error) {
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// SSE 格式: data: {...}
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))

			// 流结束标记
			if string(data) == "[DONE]" {
				return nil, io.EOF
			}

			return data, nil
		}
	}
}

// Close 关闭流
func (s *standardStreamReader) Close() error {
	return s.resp.Body.Close()
}

// APIError API 错误（为了避免循环导入，这里重新定义）
type APIError struct {
	StatusCode int
	Type       string
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s (code: %s)", e.StatusCode, e.Type, e.Message, e.Code)
}

// RateLimitError 限流错误
type RateLimitError struct {
	*APIError
	RetryAfter int
}
