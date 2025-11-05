package zhinao

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbeddings_Mock(t *testing.T) {
	// 创建 mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求
		if r.URL.Path != "/embeddings" {
			t.Errorf("Expected path /v1/embeddings, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// 返回 mock 响应
		resp := EmbeddingsResponse{
			Data: []EmbeddingData{
				{
					Embedding: []float64{0.1, 0.2, 0.3},
					Object:    "",
					Index:     0,
				},
			},
			Model:  "embedding_s1_v1",
			Object: "",
			Usage: EmbeddingsUsage{
				PromptTokens: 2,
				TotalTokens:  2,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// 创建客户端
	client, err := NewClient("test-key", WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 测试向量生成
	req := &EmbeddingsRequest{
		Model: ModelEmbeddingS1V1,
		Input: []string{"你好"},
	}

	resp, err := client.CreateEmbeddings(context.Background(), req)
	if err != nil {
		t.Fatalf("Embeddings.Create failed: %v", err)
	}

	// 验证响应
	if resp.Model != "embedding_s1_v1" {
		t.Errorf("Expected model embedding_s1_v1, got %s", resp.Model)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 embedding, got %d", len(resp.Data))
	}

	if len(resp.Data[0].Embedding) != 3 {
		t.Errorf("Expected embedding length 3, got %d", len(resp.Data[0].Embedding))
	}

	if resp.Usage.PromptTokens != 2 {
		t.Errorf("Expected prompt tokens 2, got %d", resp.Usage.PromptTokens)
	}
}

func TestEmbeddingsBuilder(t *testing.T) {
	tests := []struct {
		name     string
		build    func() *EmbeddingsRequest
		wantErr  bool
		validate func(*testing.T, *EmbeddingsRequest)
	}{
		{
			name: "basic_build",
			build: func() *EmbeddingsRequest {
				return NewEmbeddings(ModelEmbeddingS1V1).
					AddInput("你好").
					Build()
			},
			validate: func(t *testing.T, req *EmbeddingsRequest) {
				if req.Model != ModelEmbeddingS1V1 {
					t.Errorf("Expected model %s, got %s", ModelEmbeddingS1V1, req.Model)
				}
				if len(req.Input) != 1 {
					t.Errorf("Expected 1 input, got %d", len(req.Input))
				}
				if req.Input[0] != "你好" {
					t.Errorf("Expected input '你好', got '%s'", req.Input[0])
				}
			},
		},
		{
			name: "multiple_inputs",
			build: func() *EmbeddingsRequest {
				return NewEmbeddings(ModelEmbeddingS1V1).
					AddInput("你好").
					AddInput("世界").
					Build()
			},
			validate: func(t *testing.T, req *EmbeddingsRequest) {
				if len(req.Input) != 2 {
					t.Errorf("Expected 2 inputs, got %d", len(req.Input))
				}
			},
		},
		{
			name: "batch_inputs",
			build: func() *EmbeddingsRequest {
				return NewEmbeddings(ModelEmbeddingS1V1).
					AddInputs([]string{"你好", "世界", "测试"}).
					Build()
			},
			validate: func(t *testing.T, req *EmbeddingsRequest) {
				if len(req.Input) != 3 {
					t.Errorf("Expected 3 inputs, got %d", len(req.Input))
				}
			},
		},
		{
			name: "with_user",
			build: func() *EmbeddingsRequest {
				return NewEmbeddings(ModelEmbeddingS1V1).
					AddInput("你好").
					SetUser("test-user").
					Build()
			},
			validate: func(t *testing.T, req *EmbeddingsRequest) {
				if req.User != "test-user" {
					t.Errorf("Expected user 'test-user', got '%s'", req.User)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.build()
			if tt.validate != nil {
				tt.validate(t, req)
			}

			// 验证请求
			err := req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEmbeddingsRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *EmbeddingsRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid_request",
			req: &EmbeddingsRequest{
				Model: ModelEmbeddingS1V1,
				Input: []string{"你好"},
			},
			wantErr: false,
		},
		{
			name: "empty_model",
			req: &EmbeddingsRequest{
				Model: "",
				Input: []string{"你好"},
			},
			wantErr: true,
			errMsg:  "invalid model",
		},
		{
			name: "empty_input",
			req: &EmbeddingsRequest{
				Model: ModelEmbeddingS1V1,
				Input: []string{},
			},
			wantErr: true,
		},
		{
			name: "nil_input",
			req: &EmbeddingsRequest{
				Model: ModelEmbeddingS1V1,
				Input: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Logf("Got error: %v", err)
			}
		})
	}
}
