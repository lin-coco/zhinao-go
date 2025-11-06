package zhinao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateEmbeddings 生成向量
//
// 根据输入的内容，生成向量表示。
//
// 示例：
//
//	req := &zhinao.EmbeddingRequest{
//	    Model: "embedding_s1_v1",
//	    Input: []string{"你好", "世界"},
//	}
//	resp, err := client.CreateEmbeddings(ctx, req)
func (c *Client) CreateEmbeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 构建请求
	httpReq, err := c.buildRequest(ctx, "POST", "/v1/embeddings", bytes.NewReader(jsonData))
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
	resp := &EmbeddingResponse{}
	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return resp, nil
}

// EmbeddingBuilder 向量请求构建器
type EmbeddingBuilder struct {
	request *EmbeddingRequest
}

// NewEmbedding 创建向量请求构建器
//
// 示例：
//
//	req := zhinao.NewEmbedding("embedding_s1_v1").
//	    AddInput("你好").
//	    AddInput("世界").
//	    Build()
func NewEmbedding(model string) *EmbeddingBuilder {
	return &EmbeddingBuilder{
		request: &EmbeddingRequest{
			Model: model,
			Input: make([]string, 0),
		},
	}
}

// AddInput 添加输入文本
func (b *EmbeddingBuilder) AddInput(text string) *EmbeddingBuilder {
	b.request.Input = append(b.request.Input, text)
	return b
}

// AddInputs 批量添加输入文本
func (b *EmbeddingBuilder) AddInputs(texts []string) *EmbeddingBuilder {
	b.request.Input = append(b.request.Input, texts...)
	return b
}

// SetUser 设置用户标识
func (b *EmbeddingBuilder) SetUser(user string) *EmbeddingBuilder {
	b.request.User = user
	return b
}

// Build 构建请求
func (b *EmbeddingBuilder) Build() *EmbeddingRequest {
	return b.request
}
