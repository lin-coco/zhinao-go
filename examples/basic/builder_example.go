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

	// 使用 Builder 模式构建请求
	req := zhinao.NewChatBuilder().
		SetModel("360gpt-turbo").
		AddSystemMessage("你是一个专业的编程助手，擅长解释技术概念").
		AddUserMessage("请解释什么是 Go 语言的 interface").
		SetTemperature(0.7).
		SetMaxTokens(500).
		Build()

	// 发送请求
	ctx := context.Background()
	resp, err := client.Chat.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 打印响应
	if len(resp.Choices) > 0 {
		fmt.Printf("回复:\n%s\n", resp.Choices[0].Message.Content)
		fmt.Printf("\n使用token: %d (输入: %d, 输出: %d)\n",
			resp.Usage.TotalTokens,
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens)
	}
}
