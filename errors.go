package zhinao

import (
	"errors"
	"fmt"
)

var (
	// 客户端错误
	ErrMissingAPIKey = errors.New("API key is required")
	ErrInvalidConfig = errors.New("invalid configuration")

	// 请求错误
	ErrEmptyMessages = errors.New("messages cannot be empty")
	ErrInvalidModel  = errors.New("invalid model")

	// 流式错误
	ErrStreamClosed = errors.New("stream is closed")
)

// APIError API 返回的错误
type APIError struct {
	StatusCode int         `json:"status_code"`
	Type       string      `json:"type"`
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s (code: %s)",
		e.StatusCode, e.Type, e.Message, e.Code)
}

// IsRetryable 判断错误是否可重试
func (e *APIError) IsRetryable() bool {
	// 5xx 服务器错误和特定的 4xx 错误可重试
	return e.StatusCode >= 500 ||
		e.StatusCode == 429 || // Too Many Requests
		e.StatusCode == 408 // Request Timeout
}

// RateLimitError 限流错误
type RateLimitError struct {
	*APIError
	RetryAfter int // 秒数
}

func (r *RateLimitError) Error() string {
	if r.RetryAfter > 0 {
		return fmt.Sprintf("%s (retry after %d seconds)", r.APIError.Error(), r.RetryAfter)
	}
	return r.APIError.Error()
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}
