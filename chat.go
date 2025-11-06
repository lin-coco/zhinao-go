package zhinao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateChatCompletion 创建聊天补全（非流式）
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - req: 聊天请求参数
//
// 返回:
//   - *ChatCompletionResponse: 聊天响应
//   - error: 如果请求失败则返回错误
//
// 示例:
//
//	req := &zhinao.ChatCompletionRequest{
//	    Model: "360gpt-turbo",
//	    Messages: []zhinao.ChatCompletionMessage{
//	        {Role: "user", Content: "你好"},
//	    },
//	}
//	resp, err := client.CreateChatCompletion(ctx, req)
func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 确保不是流式请求
	req.Stream = false

	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 构建请求
	httpReq, err := c.buildRequest(ctx, "POST", "/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// 发送请求
	httpResp, err := c.doRequest(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// 解析响应
	var resp ChatCompletionResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// CreateChatCompletionStream 创建流式聊天补全
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - req: 聊天请求参数
//
// 返回:
//   - ChatStream: 聊天流，用于接收流式响应
//   - error: 如果请求失败则返回错误
//
// 示例:
//
//	req := &zhinao.ChatCompletionRequest{
//	    Model: "360gpt-turbo",
//	    Messages: []zhinao.ChatCompletionMessage{
//	        {Role: "user", Content: "写一首诗"},
//	    },
//	}
//	stream, err := client.CreateChatCompletionStream(ctx, req)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer stream.Close()
//
//	for {
//	    resp, err := stream.Recv()
//	    if err == io.EOF {
//	        break
//	    }
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    // 处理响应
//	}
func (c *Client) CreateChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (ChatStream, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 设置为流式
	req.Stream = true

	// 创建流
	stream, err := newChatStream(ctx, c, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
