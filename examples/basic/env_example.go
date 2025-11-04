package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	// 方法1: 使用 NewClientFromEnv（推荐）
	// 确保设置了环境变量: export ZHINAO_API_KEY=your-api-key
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 方法2: 使用 NewClient("")，效果相同
	// client, err := zhinao.NewClient("")

	// 方法3: 手动从环境变量获取
	// apiKey := os.Getenv("ZHINAO_API_KEY")
	// client, err := zhinao.NewClient(apiKey)

	// 验证环境变量
	if os.Getenv(zhinao.EnvAPIKey) == "" {
		log.Fatal("请设置环境变量 ZHINAO_API_KEY")
	}

	// 创建聊天请求
	req := &zhinao.ChatRequest{
		Model: "360gpt-turbo",
		Messages: []zhinao.Message{
			{
				Role:    "user",
				Content: "使用环境变量配置的示例",
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

	fmt.Println("\n提示: 可以通过以下方式设置环境变量")
	fmt.Println("  Linux/macOS: export ZHINAO_API_KEY=your-api-key")
	fmt.Println("  Windows:     set ZHINAO_API_KEY=your-api-key")
	fmt.Println("  或在 .env 文件中配置")
}
