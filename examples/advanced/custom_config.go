package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	// 使用自定义配置创建客户端
	client, err := zhinao.NewClient(
		"your-api-key-here",
		// 设置自定义超时时间
		zhinao.WithTimeout(30*time.Second),
		// 设置重试配置：最多重试5次，初始延迟2秒（使用指数退避）
		zhinao.WithRetry(5, 2*time.Second),
		// 设置自定义 User-Agent
		zhinao.WithUserAgent("MyApp/1.0.0"),
		// 如果需要，可以设置自定义基础URL
		// zhinao.WithBaseURL("https://custom-api.360.cn/v1"),
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 打印配置信息
	config := client.GetConfig()
	fmt.Printf("配置信息:\n")
	fmt.Printf("  - 基础URL: %s\n", config.BaseURL)
	fmt.Printf("  - 超时时间: %v\n", config.Timeout)
	fmt.Printf("  - 最大重试次数: %d\n", config.MaxRetries)
	fmt.Printf("  - 重试延迟: %v\n", config.RetryDelay)
	fmt.Println()

	// 创建带超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 发送请求
	req := &zhinao.ChatRequest{
		Model: "360gpt-turbo",
		Messages: []zhinao.Message{
			{
				Role:    "user",
				Content: "请简要介绍Go语言的并发模型",
			},
		},
	}

	resp, err := client.Chat.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	if len(resp.Choices) > 0 {
		fmt.Printf("回复: %s\n", resp.Choices[0].Message.Content)
	}
}
