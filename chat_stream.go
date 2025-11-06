package zhinao

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	stdhttp "net/http"
)

// ChatStream 聊天流接口
type ChatStream interface {
	// Recv 接收下一个响应片段
	Recv() (*ChatCompletionStreamResponse, error)
	// Close 关闭流
	Close() error
}

// StreamReader 流读取器接口
type StreamReader interface {
	// Read 读取下一行数据
	Read() ([]byte, error)
	// Close 关闭流
	Close() error
}

// streamReader 流读取器实现
type streamReader struct {
	resp   *stdhttp.Response
	reader *bufio.Reader
	ctx    context.Context    // Context for cancellation
	cancel context.CancelFunc // Cancel function for the context
}

// Read 读取下一行数据
func (s *streamReader) Read() ([]byte, error) {
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// SSE 格式: data: {...}
		if bytes.HasPrefix(line, []byte("data: ")) {
			data := bytes.TrimPrefix(line, []byte("data: "))

			// 流结束标记
			if string(data) == "[DONE]" {
				return nil, io.EOF
			}

			return data, nil
		}
	}
}

// Close 关闭流
func (s *streamReader) Close() error {
	if s.cancel != nil {
		s.cancel()
	}
	return s.resp.Body.Close()
}

// chatStream 聊天流实现
type chatStream struct {
	ctx      context.Context
	reader   StreamReader
	isClosed bool
}

// newChatStream 创建新的聊天流
func newChatStream(ctx context.Context, client *Client, req *ChatCompletionRequest) (ChatStream, error) {
	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 创建请求
	httpReq, err := client.buildRequest(ctx, "POST", "/v1/chat/completions", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	// 发送请求（流式请求不使用重试）
	resp, err := client.httpDoer.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API error (status %d): failed to read error response", resp.StatusCode)
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	// 返回流读取器
	reader := &streamReader{
		resp:   resp,
		reader: bufio.NewReader(resp.Body),
	}

	return &chatStream{
		ctx:      ctx,
		reader:   reader,
		isClosed: false,
	}, nil
}

// Recv 接收下一个响应片段
func (s *chatStream) Recv() (*ChatCompletionStreamResponse, error) {
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
	var resp ChatCompletionStreamResponse
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
