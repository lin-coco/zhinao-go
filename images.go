package zhinao

import (
	"context"
)

// ImagesService 图像服务接口
type ImagesService interface {
	// Text2Img 根据文本生成图像
	Text2Img(ctx context.Context, req *Text2ImgRequest) (*Text2ImgResponse, error)
}

// imagesService 图像服务实现
type imagesService struct {
	client *Client
}

// newImagesService 创建图像服务实例
func newImagesService(client *Client) ImagesService {
	return &imagesService{
		client: client,
	}
}

// ImageStyle 图像风格类型
type ImageStyle string

const (
	// ImageStylePapercut 剪纸风格
	ImageStylePapercut ImageStyle = "papercut"

	// ImageStyleRealistic 写实风格
	ImageStyleRealistic ImageStyle = "realistic"

	// ImageStyleCartoon 卡通风格
	ImageStyleCartoon ImageStyle = "cartoon"

	// ImageStyleCG CG风格
	ImageStyleCG ImageStyle = "CG"
)

// Text2ImgRequest 文本生成图像请求
type Text2ImgRequest struct {
	// Model 模型类型（必填）
	Model string `json:"model"`

	// Style 生成的图像风格（必填）
	// 可选风格：papercut(剪纸)、realistic(写实)、cartoon(卡通)、CG
	Style ImageStyle `json:"style"`

	// Prompt 用于生成图像的提示信息（必填）
	// 长度应小于等于 500 个字符
	Prompt string `json:"prompt"`

	// NegativePrompt 用于生成图像的负向提示信息（可选）
	// 即不希望出现的元素，长度应小于等于 500 个字符
	NegativePrompt string `json:"negative_prompt,omitempty"`

	// Width 生成的图像宽度（可选）
	// 取值应大于等于 8 小于等于 2048，默认是 512
	Width int `json:"width,omitempty"`

	// Height 生成的图像高度（可选）
	// 取值应大于等于 8 小于等于 2048，默认是 512
	Height int `json:"height,omitempty"`

	// Samples 需要生成的图片数量（可选）
	// 取值应大于等于 1 小于等于 4，默认是 1
	Samples int `json:"samples,omitempty"`

	// NumInferenceSteps 采样步数（可选）
	// 取值应大于等于 20 小于等于 50，默认是 25
	NumInferenceSteps int `json:"num_inference_steps,omitempty"`

	// Seed 种子（可选）
	// 正整数，默认随机，用于控制生成图片的随机性
	Seed int `json:"seed,omitempty"`

	// GuidanceScale 提示词强度（可选）
	// 取值应大于等于 7.5 小于等于 15，默认是 7.5
	GuidanceScale float64 `json:"guidance_scale,omitempty"`

	// EnhancePrompt 是否进行 prompt 润色（可选）
	// 默认是 false
	EnhancePrompt bool `json:"enhance_prompt,omitempty"`
}

// Text2ImgResponse 文本生成图像响应
type Text2ImgResponse struct {
	// Status 状态
	Status string `json:"status"`

	// GenerationTime 耗时，以秒为单位
	GenerationTime int `json:"generationTime"`

	// Output 生成的图片链接数组
	// 如果生成数少于指定的数，一般是因为部分图片涉及敏感被过滤
	Output []string `json:"output"`

	// Meta 元数据
	Meta ImageMeta `json:"meta"`
}

// ImageMeta 图像元数据
type ImageMeta struct {
	// H 图像高度
	H int `json:"H"`

	// W 图像宽度
	W int `json:"W"`

	// GuidanceScale 提示词强度
	GuidanceScale float64 `json:"guidance_scale"`

	// NSamples 一次请求生成的图片数
	NSamples int `json:"n_samples"`

	// NegativePrompt 用于图像生成的负向 prompt
	NegativePrompt string `json:"negative_prompt"`

	// Prompt 用于图像生成的 prompt
	Prompt string `json:"prompt"`

	// Seed 种子
	Seed int `json:"seed"`

	// Steps 采样步数
	Steps int `json:"steps"`
}

// Validate 验证请求参数
func (r *Text2ImgRequest) Validate() error {
	if r.Model == "" {
		return ErrInvalidModel
	}

	if r.Style == "" {
		return &ValidationError{Message: "style cannot be empty"}
	}

	if r.Prompt == "" {
		return &ValidationError{Message: "prompt cannot be empty"}
	}

	if len(r.Prompt) > 500 {
		return &ValidationError{Message: "prompt length must be <= 500"}
	}

	if r.NegativePrompt != "" && len(r.NegativePrompt) > 500 {
		return &ValidationError{Message: "negative_prompt length must be <= 500"}
	}

	if r.Width != 0 && (r.Width < 8 || r.Width > 2048) {
		return &ValidationError{Message: "width must be between 8 and 2048"}
	}

	if r.Height != 0 && (r.Height < 8 || r.Height > 2048) {
		return &ValidationError{Message: "height must be between 8 and 2048"}
	}

	if r.Samples != 0 && (r.Samples < 1 || r.Samples > 4) {
		return &ValidationError{Message: "samples must be between 1 and 4"}
	}

	if r.NumInferenceSteps != 0 && (r.NumInferenceSteps < 20 || r.NumInferenceSteps > 50) {
		return &ValidationError{Message: "num_inference_steps must be between 20 and 50"}
	}

	if r.GuidanceScale != 0 && (r.GuidanceScale < 7.5 || r.GuidanceScale > 15) {
		return &ValidationError{Message: "guidance_scale must be between 7.5 and 15"}
	}

	return nil
}

// Text2Img 根据文本生成图像
//
// 参数:
//   - ctx: 上下文
//   - req: 图像生成请求
//
// 返回:
//   - *Text2ImgResponse: 图像生成响应
//   - error: 错误信息
//
// 示例:
//
//	req := &zhinao.Text2ImgRequest{
//	    Model: "360CV_S0_V5",
//	    Style: zhinao.ImageStyleRealistic,
//	    Prompt: "画一个蓝天白云的图片",
//	    Width: 512,
//	    Height: 512,
//	}
//	resp, err := client.Images.Text2Img(ctx, req)
func (s *imagesService) Text2Img(ctx context.Context, req *Text2ImgRequest) (*Text2ImgResponse, error) {
	// 验证请求参数
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 发送请求
	var resp Text2ImgResponse
	err := s.client.httpClient.Post(
		ctx,
		"/images/text2img",
		req,
		&resp,
		s.client.config.APIKey,
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
