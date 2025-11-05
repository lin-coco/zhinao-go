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

	resp, err := client.CreateChatCompletion(
		context.Background(),
		&zhinao.ChatCompletionRequest{
			Model: zhinao.Model360GPTTurbo,
			Messages: []zhinao.ChatCompletionMessage{
				{
					Role:    zhinao.RoleUser,
					Content: "你好，请介绍一下360智脑",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
