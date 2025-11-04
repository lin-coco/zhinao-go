package zhinao

import (
	"os"

	"github.com/lin-coco/zhinao-go/internal/http"
)

const (
	// EnvAPIKey 环境变量名称
	EnvAPIKey = "ZHINAO_API_KEY"
)

// Client 360智脑 API 客户端
type Client struct {
	config     *Config
	httpClient http.Client

	// 服务
	Chat   ChatService
	Models ModelsService
	Images ImagesService
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
//	    zhinao.WithRetry(5, 2*time.Second),
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
		APIKey:     apiKey,
		BaseURL:    DefaultBaseURL,
		Timeout:    DefaultTimeout,
		MaxRetries: DefaultMaxRetries,
		RetryDelay: DefaultRetryDelay,
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

	// 创建 HTTP 客户端
	var httpClient http.Client
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	} else {
		httpClient = http.NewStandardClient(
			config.BaseURL,
			config.Timeout,
			config.MaxRetries,
			config.RetryDelay,
		)
	}

	client := &Client{
		config:     config,
		httpClient: httpClient,
	}

	// 初始化服务
	client.Chat = newChatService(client)
	client.Models = newModelsService(client)
	client.Images = newImagesService(client)

	return client, nil
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
