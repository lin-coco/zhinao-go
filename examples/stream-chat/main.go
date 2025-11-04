package main

import (
	"context"
	"fmt"
	"github.com/lin-coco/zhinao-go"
	"io"
)

func main() {
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		fmt.Printf("NewClientFromEnv error: %v\n", err)
		return
	}

	req := &zhinao.ChatRequest{
		Model: zhinao.Model360GPTTurbo,
		Messages: []zhinao.Message{
			{
				Role:    zhinao.RoleUser,
				Content: "请写一首关于春天的诗",
			},
		},
		Stream: true,
	}

	stream, err := client.Chat.CreateCompletionStream(context.Background(), req)
	if err != nil {
		fmt.Printf("CreateCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Println("流式响应:")
	fmt.Println("---------------------")

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("\n\n流式响应完成")
			break
		}
		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		if len(resp.Choices) > 0 {
			fmt.Print(resp.Choices[0].Delta.Content)
		}
	}
}
