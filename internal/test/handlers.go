package test

import (
	"encoding/json"
	"io"
	"net/http"
)

// ChatCompletionResponse 模拟的聊天补全响应
func ChatCompletionResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 读取请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// 解析请求
		var req map[string]interface{}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 构造模拟响应
		response := map[string]interface{}{
			"id":      "chatcmpl-test-123",
			"object":  "chat.completion",
			"created": 1234567890,
			"model":   req["model"],
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "这是一个测试响应。Hello! How can I help you today?",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 20,
				"total_tokens":      30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ChatCompletionStreamResponse 模拟的流式聊天响应
func ChatCompletionStreamResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 发送多个流式响应块
		chunks := []string{
			`data: {"id":"chatcmpl-test-123","object":"chat.completion.chunk","created":1234567890,"model":"360gpt-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}` + "\n\n",
			`data: {"id":"chatcmpl-test-123","object":"chat.completion.chunk","created":1234567890,"model":"360gpt-turbo","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}` + "\n\n",
			`data: {"id":"chatcmpl-test-123","object":"chat.completion.chunk","created":1234567890,"model":"360gpt-turbo","choices":[{"index":0,"delta":{"content":"!"},"finish_reason":null}]}` + "\n\n",
			`data: {"id":"chatcmpl-test-123","object":"chat.completion.chunk","created":1234567890,"model":"360gpt-turbo","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}` + "\n\n",
			`data: [DONE]` + "\n\n",
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		for _, chunk := range chunks {
			w.Write([]byte(chunk))
			flusher.Flush()
		}
	}
}

// ErrorResponse 模拟错误响应
func ErrorResponse(statusCode int, message, errorType, code string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"message": message,
				"type":    errorType,
				"code":    code,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}
}

// RateLimitResponse 模拟限流响应
func RateLimitResponse(retryAfter int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Rate limit exceeded",
				"type":    "rate_limit_error",
				"code":    "rate_limit",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", string(rune(retryAfter)))
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(response)
	}
}

// ModelsListResponse 模拟模型列表响应
func ModelsListResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"id":       "360gpt-turbo",
					"object":   "model",
					"created":  1234567890,
					"owned_by": "360",
				},
				{
					"id":       "360gpt-pro",
					"object":   "model",
					"created":  1234567890,
					"owned_by": "360",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
