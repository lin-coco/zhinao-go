package main

import (
	"context"
	"fmt"
	"github.com/lin-coco/zhinao-go"
)

func main() {
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		fmt.Printf("NewClientFromEnv error: %v\n", err)
		return
	}

	// 使用 Builder 模式构建请求
	req := zhinao.NewChatCompletionBuilder().
		SetModel(zhinao.Model360GPTTurbo).
		AddSystemMessage("你是一个专业的技术顾问，擅长解释复杂的技术概念").
		AddUserMessage("请用简单的语言解释什么是机器学习").
		SetTemperature(0.7).
		SetMaxTokens(500).
		SetTopP(0.9).
		Build()

	fmt.Println("使用 Builder 模式构建的请求:")
	fmt.Printf("- Model: %s\n", req.Model)
	fmt.Printf("- Temperature: %.1f\n", req.Temperature)
	fmt.Printf("- MaxTokens: %d\n", req.MaxTokens)
	fmt.Printf("- TopP: %.1f\n", req.TopP)
	fmt.Printf("- Messages: %d 条\n\n", len(req.Messages))

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println("AI 回复:")
	fmt.Println("---------------------")
	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Printf("\n使用 Token 数: %d\n", resp.Usage.TotalTokens)
}
