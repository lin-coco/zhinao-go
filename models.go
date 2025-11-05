package zhinao

import (
	"context"
	"encoding/json"
	"fmt"
)

// ListModels 获取可用模型列表
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
//	resp, err := client.ListModels(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, model := range resp.Data {
//	    fmt.Printf("Model: %s\n", model.ID)
//	}
func (c *Client) ListModels(ctx context.Context) (*ModelsResponse, error) {
	// 构建请求
	httpReq, err := c.buildRequest(ctx, "GET", "/models", nil)
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
	var resp ModelsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// GetModel 获取指定模型信息
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
//	info, err := client.GetModel(ctx, "360gpt-turbo")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Model: %s, Owner: %s\n", info.ID, info.OwnedBy)
func (c *Client) GetModel(ctx context.Context, modelID string) (*ModelInfo, error) {
	if modelID == "" {
		return nil, fmt.Errorf("model ID cannot be empty")
	}

	// 构建请求
	httpReq, err := c.buildRequest(ctx, "GET", "/models/"+modelID, nil)
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
	var resp ModelInfo
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}
