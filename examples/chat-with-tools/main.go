package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lin-coco/zhinao-go"
)

func main() {
	client, err := zhinao.NewClientFromEnv()
	if err != nil {
		fmt.Printf("NewClientFromEnv error: %v\n", err)
		return
	}

	// 定义天气查询工具
	weatherTool := zhinao.Tool{
		Type: "function",
		Function: zhinao.ToolFunction{
			Name:        "get_current_weather",
			Description: "获取指定城市的当前天气",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "城市名称，例如：北京、上海",
					},
					"unit": map[string]interface{}{
						"type": "string",
						"enum": []string{"celsius", "fahrenheit"},
					},
				},
				"required": []string{"location"},
			},
		},
	}

	// 初始对话
	messages := []zhinao.Message{
		{
			Role:    zhinao.RoleUser,
			Content: "北京今天天气怎么样？",
		},
	}

	fmt.Println("询问 AI: '北京今天天气怎么样？'")
	fmt.Println("提供工具: get_current_weather()")

	// 第一次调用：AI 决定是否使用工具
	resp, err := client.CreateCompletion(
		context.Background(),
		&zhinao.ChatRequest{
			Model:    zhinao.Model360GPTTurbo,
			Messages: messages,
			Tools:    []zhinao.Tool{weatherTool},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	msg := resp.Choices[0].Message

	// 检查 AI 是否请求调用工具
	if len(msg.ToolCalls) == 0 {
		fmt.Printf("\nAI 直接回复: %s\n", msg.Content)
		return
	}

	// AI 请求调用工具
	toolCall := msg.ToolCalls[0]
	fmt.Printf("\nAI 请求调用工具: %s\n", toolCall.Function.Name)
	fmt.Printf("参数: %s\n", toolCall.Function.Arguments)

	// 解析参数
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &params); err != nil {
		fmt.Printf("Parse arguments error: %v\n", err)
		return
	}

	// 模拟调用天气 API（实际应用中这里应该调用真实的天气 API）
	weatherResult := map[string]interface{}{
		"location":    params["location"],
		"temperature": 22,
		"unit":        "celsius",
		"description": "晴朗",
	}
	weatherJSON, _ := json.Marshal(weatherResult)

	fmt.Printf("\n模拟工具返回: %s\n", string(weatherJSON))

	// 将工具调用和结果添加到对话历史
	messages = append(messages, msg)
	messages = append(messages, zhinao.Message{
		Role:       zhinao.RoleTool,
		Content:    string(weatherJSON),
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
	})

	// 第二次调用：让 AI 根据工具结果给出最终答案
	fmt.Println("\n请求 AI 基于工具结果回答原始问题...")
	resp, err = client.CreateCompletion(
		context.Background(),
		&zhinao.ChatRequest{
			Model:    zhinao.Model360GPTTurbo,
			Messages: messages,
			Tools:    []zhinao.Tool{weatherTool},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Printf("\nAI 最终回复: %s\n", resp.Choices[0].Message.Content)
}
