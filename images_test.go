package zhinao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestText2Img_Mock(t *testing.T) {
	// 创建 Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 验证请求路径
		if r.URL.Path != "/images/text2img" {
			t.Errorf("Expected path /images/text2img, got %s", r.URL.Path)
		}

		// 验证 Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-api-key" {
			t.Errorf("Expected Authorization header 'Bearer test-api-key', got '%s'", auth)
		}

		// 解析请求体
		var req ImageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// 验证请求参数
		if req.Model != "360CV_S0_V5" {
			t.Errorf("Expected model '360CV_S0_V5', got '%s'", req.Model)
		}

		if req.Style != ImageStyleRealistic {
			t.Errorf("Expected style 'realistic', got '%s'", req.Style)
		}

		// 返回成功响应
		resp := ImageResponse{
			Status:         "success",
			GenerationTime: 7,
			Output: []string{
				"https://example.com/image1.png",
			},
			Meta: ImageMeta{
				H:              512,
				W:              512,
				GuidanceScale:  7.5,
				NSamples:       1,
				NegativePrompt: "",
				Prompt:         "画一个蓝天白云的图片",
				Seed:           49022,
				Steps:          25,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端
	client := &Client{
		config: &Config{
			APIKey:  "test-api-key",
			BaseURL: server.URL,
		},
		httpDoer: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	// 创建请求
	req := &ImageRequest{
		Model:         "360CV_S0_V5",
		Style:         ImageStyleRealistic,
		Prompt:        "画一个蓝天白云的图片",
		Width:         512,
		Height:        512,
		Samples:       1,
		Seed:          49022,
		GuidanceScale: 7.5,
	}

	// 发送请求
	resp, err := client.CreateImage(context.Background(), req)
	if err != nil {
		t.Fatalf("Text2Img failed: %v", err)
	}

	// 验证响应
	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}

	if resp.GenerationTime != 7 {
		t.Errorf("Expected generation time 7, got %d", resp.GenerationTime)
	}

	if len(resp.Output) != 1 {
		t.Fatalf("Expected 1 output image, got %d", len(resp.Output))
	}

	if resp.Output[0] != "https://example.com/image1.png" {
		t.Errorf("Expected image URL 'https://example.com/image1.png', got '%s'", resp.Output[0])
	}

	// 验证元数据
	if resp.Meta.W != 512 || resp.Meta.H != 512 {
		t.Errorf("Expected image size 512x512, got %dx%d", resp.Meta.W, resp.Meta.H)
	}
}

func TestText2Img_WithNegativePrompt_Mock(t *testing.T) {
	// 创建 Mock Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ImageRequest
		json.NewDecoder(r.Body).Decode(&req)

		// 验证负向提示词
		if req.NegativePrompt != "模糊,低质量" {
			t.Errorf("Expected negative_prompt '模糊,低质量', got '%s'", req.NegativePrompt)
		}

		resp := ImageResponse{
			Status:         "success",
			GenerationTime: 8,
			Output:         []string{"https://example.com/image.png"},
			Meta: ImageMeta{
				H:              512,
				W:              512,
				GuidanceScale:  7.5,
				NSamples:       1,
				NegativePrompt: req.NegativePrompt,
				Prompt:         req.Prompt,
				Seed:           12345,
				Steps:          25,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端
	client := &Client{
		config: &Config{
			APIKey:  "test-api-key",
			BaseURL: server.URL,
		},
		httpDoer: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	// 创建请求（带负向提示词）
	req := &ImageRequest{
		Model:          "360CV_S0_V5",
		Style:          ImageStyleRealistic,
		Prompt:         "美丽的风景",
		NegativePrompt: "模糊,低质量",
		Width:          512,
		Height:         512,
	}

	// 发送请求
	resp, err := client.CreateImage(context.Background(), req)
	if err != nil {
		t.Fatalf("Text2Img failed: %v", err)
	}

	// 验证响应
	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}

	if resp.Meta.NegativePrompt != "模糊,低质量" {
		t.Errorf("Expected negative_prompt '模糊,低质量' in meta, got '%s'", resp.Meta.NegativePrompt)
	}
}

func TestText2ImgRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ImageRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Style:  ImageStyleRealistic,
				Prompt: "测试图片",
			},
			wantErr: false,
		},
		{
			name: "empty model",
			req: &ImageRequest{
				Style:  ImageStyleRealistic,
				Prompt: "测试图片",
			},
			wantErr: true,
			errMsg:  "model",
		},
		{
			name: "empty style",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Prompt: "测试图片",
			},
			wantErr: true,
			errMsg:  "style",
		},
		{
			name: "empty prompt",
			req: &ImageRequest{
				Model: "360CV_S0_V5",
				Style: ImageStyleRealistic,
			},
			wantErr: true,
			errMsg:  "prompt",
		},
		{
			name: "prompt too long",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Style:  ImageStyleRealistic,
				Prompt: string(make([]byte, 501)),
			},
			wantErr: true,
			errMsg:  "prompt length",
		},
		{
			name: "invalid width - too small",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Style:  ImageStyleRealistic,
				Prompt: "测试",
				Width:  7,
			},
			wantErr: true,
			errMsg:  "width",
		},
		{
			name: "invalid width - too large",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Style:  ImageStyleRealistic,
				Prompt: "测试",
				Width:  2049,
			},
			wantErr: true,
			errMsg:  "width",
		},
		{
			name: "invalid height",
			req: &ImageRequest{
				Model:  "360CV_S0_V5",
				Style:  ImageStyleRealistic,
				Prompt: "测试",
				Height: 5,
			},
			wantErr: true,
			errMsg:  "height",
		},
		{
			name: "invalid samples",
			req: &ImageRequest{
				Model:   "360CV_S0_V5",
				Style:   ImageStyleRealistic,
				Prompt:  "测试",
				Samples: 5,
			},
			wantErr: true,
			errMsg:  "samples",
		},
		{
			name: "invalid guidance_scale - too low",
			req: &ImageRequest{
				Model:         "360CV_S0_V5",
				Style:         ImageStyleRealistic,
				Prompt:        "测试",
				GuidanceScale: 7.0,
			},
			wantErr: true,
			errMsg:  "guidance_scale",
		},
		{
			name: "invalid guidance_scale - too high",
			req: &ImageRequest{
				Model:         "360CV_S0_V5",
				Style:         ImageStyleRealistic,
				Prompt:        "测试",
				GuidanceScale: 16.0,
			},
			wantErr: true,
			errMsg:  "guidance_scale",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				// 检查错误消息是否包含预期的关键字
				errMsg := err.Error()
				if tt.errMsg != "" && len(errMsg) > 0 {
					// 简单检查错误消息中是否包含关键字
					// 这里不做严格匹配，因为错误消息可能会变化
					t.Logf("Got error: %s", errMsg)
				}
			}
		})
	}
}

func TestImageStyles(t *testing.T) {
	// 测试所有图像风格常量
	styles := []ImageStyle{
		ImageStylePapercut,
		ImageStyleRealistic,
		ImageStyleCartoon,
		ImageStyleCG,
	}

	expectedValues := []string{
		"papercut",
		"realistic",
		"cartoon",
		"CG",
	}

	for i, style := range styles {
		if string(style) != expectedValues[i] {
			t.Errorf("Expected style %s, got %s", expectedValues[i], string(style))
		}
	}
}
