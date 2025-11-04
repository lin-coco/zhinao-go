package zhinao

import (
	"context"
)

// ChatService 聊天服务接口
type ChatService interface {
	// CreateCompletion 创建聊天补全（非流式）
	CreateCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// CreateCompletionStream 创建流式聊天补全
	CreateCompletionStream(ctx context.Context, req *ChatRequest) (ChatStream, error)
}

// ChatStream 聊天流接口
type ChatStream interface {
	// Recv 接收下一个响应片段
	Recv() (*ChatStreamResponse, error)

	// Close 关闭流
	Close() error
}

// chatService 聊天服务实现
type chatService struct {
	client *Client
}

// newChatService 创建聊天服务实例
func newChatService(client *Client) ChatService {
	return &chatService{client: client}
}

// CreateCompletion 创建聊天补全（非流式）
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - req: 聊天请求参数
//
// 返回:
//   - *ChatResponse: 聊天响应
//   - error: 如果请求失败则返回错误
//
// 示例:
//
//	req := &zhinao.ChatRequest{
//	    Model: "360gpt-turbo",
//	    Messages: []zhinao.Message{
//	        {Role: "user", Content: "你好"},
//	    },
//	}
//	resp, err := client.Chat.CreateCompletion(ctx, req)
func (s *chatService) CreateCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 确保不是流式请求
	req.Stream = false

	// 发送请求
	var resp ChatResponse
	err := s.client.httpClient.Post(
		ctx,
		"/chat/completions",
		req,
		&resp,
		s.client.config.APIKey,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateCompletionStream 创建流式聊天补全
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
//	req := &zhinao.ChatRequest{
//	    Model: "360gpt-turbo",
//	    Messages: []zhinao.Message{
//	        {Role: "user", Content: "写一首诗"},
//	    },
//	}
//	stream, err := client.Chat.CreateCompletionStream(ctx, req)
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
func (s *chatService) CreateCompletionStream(ctx context.Context, req *ChatRequest) (ChatStream, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 设置为流式
	req.Stream = true

	// 创建流
	stream, err := newChatStream(ctx, s.client, req)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
