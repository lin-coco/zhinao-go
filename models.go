package zhinao

import (
	"context"
	"fmt"
)

// ModelsService 模型服务接口
type ModelsService interface {
	// List 获取可用模型列表
	List(ctx context.Context) (*ModelsResponse, error)

	// Get 获取指定模型信息
	Get(ctx context.Context, modelID string) (*ModelInfo, error)
}

// modelsService 模型服务实现
type modelsService struct {
	client *Client
}

// newModelsService 创建模型服务实例
func newModelsService(client *Client) ModelsService {
	return &modelsService{client: client}
}

// List 获取可用模型列表
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//
// 返回:
//   - *ModelsResponse: 模型列表响应
//   - error: 如果请求失败则返回错误
//
// 示例:
//
//	resp, err := client.Models.List(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, model := range resp.Data {
//	    fmt.Printf("Model: %s\n", model.ID)
//	}
func (s *modelsService) List(ctx context.Context) (*ModelsResponse, error) {
	var resp ModelsResponse
	err := s.client.httpClient.Get(
		ctx,
		"/models",
		&resp,
		s.client.config.APIKey,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Get 获取指定模型信息
//
// 参数:
//   - ctx: 上下文，用于控制请求的生命周期
//   - modelID: 模型 ID
//
// 返回:
//   - *ModelInfo: 模型信息
//   - error: 如果请求失败则返回错误
//
// 示例:
//
//	info, err := client.Models.Get(ctx, "360gpt-turbo")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Model: %s, Owner: %s\n", info.ID, info.OwnedBy)
func (s *modelsService) Get(ctx context.Context, modelID string) (*ModelInfo, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID cannot be empty")
	}

	var resp ModelInfo
	err := s.client.httpClient.Get(
		ctx,
		"/models/"+modelID,
		&resp,
		s.client.config.APIKey,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
