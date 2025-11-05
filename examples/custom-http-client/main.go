package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	fmt.Println("=== 自定义 HTTP 客户端示例 ===")
	fmt.Println()

	// 示例1: 使用默认配置（推荐方式）
	// WithBaseURL、WithTimeout 会自动创建合适的 HTTP 客户端
	fmt.Println("示例1: 使用默认配置")
	client1, err := zhinao.NewClient(
		"your-api-key",
		zhinao.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 客户端创建成功（使用内置配置）")
	fmt.Println()

	// 示例2: 使用自定义的 http.Client
	// 当你需要完全控制 HTTP 客户端的行为时使用
	// 适用场景：
	// - 自定义连接池配置
	// - 设置代理
	// - 配置 TLS
	// - 添加自定义 Transport
	fmt.Println("示例2: 使用自定义 http.Client")
	customHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 10,               // 每个 host 的最大空闲连接
			IdleConnTimeout:     90 * time.Second, // 空闲连接超时
			TLSHandshakeTimeout: 10 * time.Second, // TLS 握手超时
		},
	}

	client2, err := zhinao.NewClient(
		"your-api-key",
		zhinao.WithHTTPClient(customHTTPClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 客户端创建成功（使用自定义 http.Client）")
	fmt.Println()

	// 示例3: 使用自定义实现的 HTTPDoer 接口
	// 实现 HTTPDoer 接口可以完全自定义请求行为
	// 适用场景：
	// - 添加请求/响应日志
	// - 实现自定义重试逻辑
	// - 添加请求追踪
	// - 集成第三方 HTTP 库（如 go-resty）
	fmt.Println("示例3: 使用自定义 HTTPDoer 实现")
	customDoer := &CustomHTTPDoer{
		client: &http.Client{Timeout: 30 * time.Second},
	}

	client3, err := zhinao.NewClient(
		"your-api-key",
		zhinao.WithHTTPClient(customDoer),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ 客户端创建成功（使用自定义 HTTPDoer）")
	fmt.Println()

	fmt.Println("所有客户端创建成功！")
	fmt.Println("\n提示：")
	fmt.Println("- 使用 WithHTTPClient 时，WithTimeout 配置会被忽略")
	fmt.Println("- WithBaseURL 和 WithHeaders 仍然有效")
	fmt.Println("- 可以实现 HTTPDoer 接口来适配任何 HTTP 库")

	// 使用客户端
	_ = client1
	_ = client2
	_ = client3
}

// CustomHTTPDoer 是一个自定义的 HTTP 客户端实现
// 它实现了 zhinao.HTTPDoer 接口
//
// 这个示例展示了如何：
// 1. 实现 HTTPDoer 接口
// 2. 添加请求日志
// 3. 添加自定义请求头
// 4. 添加响应日志
type CustomHTTPDoer struct {
	client *http.Client
}

// Do 实现 zhinao.HTTPDoer 接口
func (c *CustomHTTPDoer) Do(req *http.Request) (*http.Response, error) {
	// 请求前的处理
	fmt.Printf("→ 发送请求: %s %s\n", req.Method, req.URL.String())

	// 添加自定义请求头（例如追踪 ID）
	req.Header.Set("X-Custom-Header", "custom-value")
	req.Header.Set("X-Request-ID", fmt.Sprintf("%d", time.Now().UnixNano()))

	// 记录请求头（调试用）
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("  Header: %s: %s\n", key, value)
		}
	}

	// 执行实际的 HTTP 请求
	start := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("✗ 请求失败: %v (耗时: %v)\n", err, duration)
		return nil, err
	}

	// 响应后的处理
	fmt.Printf("✓ 收到响应: %d (耗时: %v)\n", resp.StatusCode, duration)

	return resp, nil
}
