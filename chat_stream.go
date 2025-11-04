package zhinao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/lin-coco/zhinao-go/internal/http"
)

// chatStream 聊天流实现
type chatStream struct {
	ctx      context.Context
	reader   http.StreamReader
	isClosed bool
}

// newChatStream 创建新的聊天流
func newChatStream(ctx context.Context, client *Client, req *ChatRequest) (ChatStream, error) {
	// 发送流式请求
	reader, err := client.httpClient.PostStream(
		ctx,
		"/chat/completions",
		req,
		client.config.APIKey,
	)
	if err != nil {
		return nil, err
	}

	return &chatStream{
		ctx:      ctx,
		reader:   reader,
		isClosed: false,
	}, nil
}

// Recv 接收下一个响应片段
func (s *chatStream) Recv() (*ChatStreamResponse, error) {
	if s.isClosed {
		return nil, ErrStreamClosed
	}

	// 检查 context 是否已取消
	select {
	case <-s.ctx.Done():
		s.Close()
		return nil, s.ctx.Err()
	default:
	}

	// 读取数据
	data, err := s.reader.Read()
	if err != nil {
		if err == io.EOF {
			s.isClosed = true
		}
		return nil, err
	}

	// 解析响应
	var resp ChatStreamResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stream response: %w", err)
	}

	return &resp, nil
}

// Close 关闭流
func (s *chatStream) Close() error {
	if s.isClosed {
		return nil
	}
	s.isClosed = true
	return s.reader.Close()
}
