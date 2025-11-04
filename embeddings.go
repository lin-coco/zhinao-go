package zhinao

import (
	"context"
)

// EmbeddingsService 向量服务接口
type EmbeddingsService interface {
	// Create 生成向量
	Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error)
}

// embeddingsService 向量服务实现
type embeddingsService struct {
	client *Client
}

// newEmbeddingsService 创建向量服务实例
func newEmbeddingsService(client *Client) EmbeddingsService {
	return &embeddingsService{
		client: client,
	}
}

// EmbeddingsRequest 向量生成请求
type EmbeddingsRequest struct {
	// Model 模型类型（必填）
	Model string `json:"model"`

	// Input 批量生成，每一项是需要生成向量的内容（必填）
	Input []string `json:"input"`

	// User 标记业务方用户 id，便于业务方区分不同用户（可选）
	User string `json:"user,omitempty"`
}

// EmbeddingsResponse 向量生成响应
type EmbeddingsResponse struct {
	// Data 返回结果，一个结构体数组，每个元素对应 input 的每一项输入
	Data []EmbeddingData `json:"data"`

	// Model 本次调用使用的模型名
	Model string `json:"model"`

	// Object 对象类型
	Object string `json:"object"`

	// Usage token 消耗量
	Usage EmbeddingsUsage `json:"usage"`
}

// EmbeddingData 向量数据
type EmbeddingData struct {
	// Embedding 返回的向量结果，是一个浮点数数组
	Embedding []float64 `json:"embedding"`

	// Object 对象类型
	Object string `json:"object"`

	// Index 代表本次结果在 data 里的下标值
	Index int `json:"index"`
}

// EmbeddingsUsage token 使用量
type EmbeddingsUsage struct {
	// PromptTokens 输入 token 消耗量
	PromptTokens int `json:"prompt_tokens"`

	// TotalTokens 总 token 消耗量
	TotalTokens int `json:"total_tokens"`
}

// Create 生成向量
//
// 根据输入的内容，生成向量表示。
//
// 示例：
//
//	req := &zhinao.EmbeddingsRequest{
//	    Model: "embedding_s1_v1",
//	    Input: []string{"你好", "世界"},
//	}
//	resp, err := client.Embeddings.Create(ctx, req)
func (s *embeddingsService) Create(ctx context.Context, req *EmbeddingsRequest) (*EmbeddingsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp := &EmbeddingsResponse{}
	err := s.client.httpClient.Post(
		ctx,
		"/embeddings",
		req,
		resp,
		s.client.config.APIKey,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Validate 验证请求参数
func (req *EmbeddingsRequest) Validate() error {
	if req.Model == "" {
		return ErrInvalidModel
	}

	if len(req.Input) == 0 {
		return ErrEmptyInput
	}

	return nil
}

// EmbeddingsBuilder 向量请求构建器
type EmbeddingsBuilder struct {
	request *EmbeddingsRequest
}

// NewEmbeddings 创建向量请求构建器
//
// 示例：
//
//	req := zhinao.NewEmbeddings("embedding_s1_v1").
//	    AddInput("你好").
//	    AddInput("世界").
//	    Build()
func NewEmbeddings(model string) *EmbeddingsBuilder {
	return &EmbeddingsBuilder{
		request: &EmbeddingsRequest{
			Model: model,
			Input: make([]string, 0),
		},
	}
}

// AddInput 添加输入文本
func (b *EmbeddingsBuilder) AddInput(text string) *EmbeddingsBuilder {
	b.request.Input = append(b.request.Input, text)
	return b
}

// AddInputs 批量添加输入文本
func (b *EmbeddingsBuilder) AddInputs(texts []string) *EmbeddingsBuilder {
	b.request.Input = append(b.request.Input, texts...)
	return b
}

// SetUser 设置用户标识
func (b *EmbeddingsBuilder) SetUser(user string) *EmbeddingsBuilder {
	b.request.User = user
	return b
}

// Build 构建请求
func (b *EmbeddingsBuilder) Build() *EmbeddingsRequest {
	return b.request
}
