package zhinao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// CreateImage 根据文本生成图像
//
// 参数:
//   - ctx: 上下文
//   - req: 图像生成请求
//
// 返回:
//   - *ImageResponse: 图像生成响应
//   - error: 错误信息
//
// 示例:
//
//	req := &zhinao.ImageRequest{
//	    Model: "360CV_S0_V5",
//	    Style: zhinao.ImageStyleRealistic,
//	    Prompt: "画一个蓝天白云的图片",
//	    Width: 512,
//	    Height: 512,
//	}
//	resp, err := client.CreateImage(ctx, req)
func (c *Client) CreateImage(ctx context.Context, req *ImageRequest) (*ImageResponse, error) {
	// 验证请求参数
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 序列化请求体
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 构建请求
	httpReq, err := c.buildRequest(ctx, "POST", "/v1/images/text2img", bytes.NewReader(jsonData))
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
	var resp ImageResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}
