package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	// 创建客户端
	client, err := zhinao.NewClient("your-api-key-here")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 创建聊天请求
	req := &zhinao.ChatRequest{
		Model: "360gpt-turbo",
		Messages: []zhinao.Message{
			{
				Role:    "user",
				Content: "你好，请介绍一下360智脑",
			},
		},
	}

	// 发送请求
	ctx := context.Background()
	resp, err := client.Chat.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 打印响应
	if len(resp.Choices) > 0 {
		fmt.Printf("回复: %s\n", resp.Choices[0].Message.Content)
		fmt.Printf("使用token: %d\n", resp.Usage.TotalTokens)
	}
}
