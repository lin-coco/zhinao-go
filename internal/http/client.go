package http

import "context"

// Client HTTP 客户端接口
type Client interface {
	// Post 发送 POST 请求
	Post(ctx context.Context, path string, body, result interface{}, apiKey string) error

	// Get 发送 GET 请求
	Get(ctx context.Context, path string, result interface{}, apiKey string) error

	// PostStream 发送流式 POST 请求
	PostStream(ctx context.Context, path string, body interface{}, apiKey string) (StreamReader, error)
}

// StreamReader 流式读取接口
type StreamReader interface {
	// Read 读取下一行数据
	Read() ([]byte, error)

	// Close 关闭流
	Close() error
}
