package zhinao

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestModelsList_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册模型列表处理器
	server.RegisterHandler("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		resp := &ModelsResponse{
			Object: "list",
			Data: []ModelInfo{
				{
					ID:      Model360GPTTurbo,
					Object:  "model",
					OwnedBy: "360",
				},
				{
					ID:      "360gpt-pro",
					Object:  "model",
					OwnedBy: "360",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	resp, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("Models.List error: %v", err)
	}

	if resp.Object != "list" {
		t.Errorf("Expected object 'list', got %s", resp.Object)
	}

	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 models, got %d", len(resp.Data))
	}

	if resp.Data[0].ID != Model360GPTTurbo {
		t.Errorf("Expected first model ID %s, got %s", Model360GPTTurbo, resp.Data[0].ID)
	}
}

func TestModelsGet_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册获取模型处理器
	server.RegisterHandler("/v1/models/360gpt-turbo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		resp := &ModelInfo{
			ID:      Model360GPTTurbo,
			Object:  "model",
			OwnedBy: "360",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	info, err := client.GetModel(context.Background(), Model360GPTTurbo)
	if err != nil {
		t.Fatalf("Models.Get error: %v", err)
	}

	if info.ID != Model360GPTTurbo {
		t.Errorf("Expected model ID %s, got %s", Model360GPTTurbo, info.ID)
	}

	if info.Object != "model" {
		t.Errorf("Expected object 'model', got %s", info.Object)
	}

	if info.OwnedBy != "360" {
		t.Errorf("Expected owned by '360', got %s", info.OwnedBy)
	}
}

func TestModelsGet_EmptyID_Mock(t *testing.T) {
	client, _, teardown := setupZhinaoTestServer()
	defer teardown()

	_, err := client.GetModel(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty model ID, got nil")
	}

	expectedMsg := "model ID cannot be empty"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestModelsGet_NotFound_Mock(t *testing.T) {
	client, server, teardown := setupZhinaoTestServer()
	defer teardown()

	// 注册 404 处理器
	server.RegisterHandler("/v1/models/non-existent-model", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		errResp := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Model not found",
				"type":    "invalid_request_error",
				"code":    "model_not_found",
			},
		}
		json.NewEncoder(w).Encode(errResp)
	})

	_, err := client.GetModel(context.Background(), "non-existent-model")
	if err == nil {
		t.Error("Expected error for non-existent model, got nil")
	}
}
