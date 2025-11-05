package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/lin-coco/zhinao-go"
)

func main() {
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		fmt.Printf("NewClientFromEnv error: %v\n", err)
		return
	}

	messages := []zhinao.ChatCompletionMessage{
		{
			Role:    zhinao.RoleSystem,
			Content: "你是一个helpful的AI助手",
		},
	}

	fmt.Println("聊天机器人 (输入 'exit' 退出)")
	fmt.Println("---------------------")
	fmt.Print("> ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userInput := scanner.Text()
		if userInput == "exit" {
			fmt.Println("再见!")
			break
		}

		// 添加用户消息
		messages = append(messages, zhinao.ChatCompletionMessage{
			Role:    zhinao.RoleUser,
			Content: userInput,
		})

		// 调用 API
		resp, err := client.CreateChatCompletion(
			context.Background(),
			&zhinao.ChatCompletionRequest{
				Model:    zhinao.Model360GPTTurbo,
				Messages: messages,
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			fmt.Print("> ")
			continue
		}

		// 显示助手回复
		assistantMsg := resp.Choices[0].Message.Content
		fmt.Printf("%s\n\n", assistantMsg)

		// 添加助手消息到历史
		messages = append(messages, zhinao.ChatCompletionMessage{
			Role:    zhinao.RoleAssistant,
			Content: assistantMsg,
		})

		fmt.Print("> ")
	}
}
