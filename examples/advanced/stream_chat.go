package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	// 创建客户端
	client, err := zhinao.NewClient("your-api-key-here")
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 创建流式请求
	req := &zhinao.ChatRequest{
		Model: "360gpt-turbo",
		Messages: []zhinao.Message{
			{
				Role:    "user",
				Content: "请写一首关于春天的诗",
			},
		},
	}

	// 创建流
	ctx := context.Background()
	stream, err := client.Chat.CreateCompletionStream(ctx, req)
	if err != nil {
		log.Fatalf("创建流失败: %v", err)
	}
	defer stream.Close()

	fmt.Println("开始接收流式响应:")
	fmt.Println("---")

	// 接收流式响应
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\n---")
			fmt.Println("流结束")
			break
		}
		if err != nil {
			log.Fatalf("接收流失败: %v", err)
		}

		// 打印增量内容
		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Delta.Content
			if content != "" {
				fmt.Print(content)
			}
		}
	}
}
