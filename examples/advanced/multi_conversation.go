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

	ctx := context.Background()

	// 使用 Builder 维护对话历史
	builder := zhinao.NewChatBuilder().
		SetModel("360gpt-turbo").
		AddSystemMessage("你是一个有帮助的助手")

	// 第一轮对话
	builder.AddUserMessage("什么是人工智能？")

	req1 := builder.Build()
	resp1, err := client.Chat.CreateCompletion(ctx, req1)
	if err != nil {
		log.Fatalf("第一轮对话失败: %v", err)
	}

	fmt.Println("用户: 什么是人工智能？")
	fmt.Printf("助手: %s\n\n", resp1.Choices[0].Message.Content)

	// 将助手的回复加入对话历史
	builder.AddAssistantMessage(resp1.Choices[0].Message.Content)

	// 第二轮对话
	builder.AddUserMessage("它有哪些应用场景？")

	req2 := builder.Build()
	resp2, err := client.Chat.CreateCompletion(ctx, req2)
	if err != nil {
		log.Fatalf("第二轮对话失败: %v", err)
	}

	fmt.Println("用户: 它有哪些应用场景？")
	fmt.Printf("助手: %s\n\n", resp2.Choices[0].Message.Content)

	// 将助手的回复加入对话历史
	builder.AddAssistantMessage(resp2.Choices[0].Message.Content)

	// 第三轮对话
	builder.AddUserMessage("请举一个具体的例子")

	req3 := builder.Build()
	resp3, err := client.Chat.CreateCompletion(ctx, req3)
	if err != nil {
		log.Fatalf("第三轮对话失败: %v", err)
	}

	fmt.Println("用户: 请举一个具体的例子")
	fmt.Printf("助手: %s\n\n", resp3.Choices[0].Message.Content)

	// 打印总token使用情况
	totalTokens := resp1.Usage.TotalTokens + resp2.Usage.TotalTokens + resp3.Usage.TotalTokens
	fmt.Printf("总计使用 token: %d\n", totalTokens)
}
