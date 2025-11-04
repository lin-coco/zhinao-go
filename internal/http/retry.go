package http

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int             // 最大重试次数
	RetryDelay time.Duration   // 基础重试延迟
	Backoff    BackoffStrategy // 退避策略
}

// BackoffStrategy 退避策略
type BackoffStrategy int

const (
	// ConstantBackoff 常量退避 - 每次重试使用相同的延迟
	ConstantBackoff BackoffStrategy = iota
	// ExponentialBackoff 指数退避 - 延迟随重试次数指数增长
	ExponentialBackoff
)

// shouldRetry 判断是否应该重试
func shouldRetry(statusCode int, err error) bool {
	// 网络错误总是重试
	if err != nil {
		return true
	}

	// 5xx 服务器错误可重试
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// 特定的 4xx 错误可重试
	switch statusCode {
	case http.StatusTooManyRequests: // 429
		return true
	case http.StatusRequestTimeout: // 408
		return true
	default:
		return false
	}
}

// calculateDelay 计算重试延迟时间
func calculateDelay(config *RetryConfig, attempt int) time.Duration {
	if config.Backoff == ExponentialBackoff {
		// 指数退避: delay * 2^attempt
		multiplier := math.Pow(2, float64(attempt))
		delay := time.Duration(float64(config.RetryDelay) * multiplier)
		// 限制最大延迟为 60 秒
		maxDelay := 60 * time.Second
		if delay > maxDelay {
			return maxDelay
		}
		return delay
	}
	// 常量退避
	return config.RetryDelay
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func() (*http.Response, error)

// executeWithRetry 执行带重试的函数
func executeWithRetry(ctx context.Context, config *RetryConfig, fn RetryableFunc) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 执行请求
		resp, lastErr = fn()

		// 第一次尝试或不需要重试
		if attempt == 0 || !shouldRetry(getStatusCode(resp), lastErr) {
			return resp, lastErr
		}

		// 如果需要重试且还有重试次数
		if attempt < config.MaxRetries {
			delay := calculateDelay(config, attempt)

			// 等待重试延迟
			select {
			case <-time.After(delay):
				// 继续重试
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return resp, fmt.Errorf("max retries (%d) exceeded: %w", config.MaxRetries, lastErr)
}

// getStatusCode 从响应中获取状态码
func getStatusCode(resp *http.Response) int {
	if resp == nil {
		return 0
	}
	return resp.StatusCode
}
