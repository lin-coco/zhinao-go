package test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
)

const testAPIKey = "test-api-key-do-not-use-in-production"

// GetTestToken 返回测试用的 API Key
func GetTestToken() string {
	return testAPIKey
}

// ServerTest 测试服务器
type ServerTest struct {
	handlers map[string]http.HandlerFunc
}

// NewTestServer 创建新的测试服务器
func NewTestServer() *ServerTest {
	return &ServerTest{
		handlers: make(map[string]http.HandlerFunc),
	}
}

// RegisterHandler 注册路径处理器
// 支持正则表达式路径，例如 "/v1/chat/.*"
func (s *ServerTest) RegisterHandler(path string, handler http.HandlerFunc) {
	// 将通配符转换为正则表达式
	path = strings.ReplaceAll(path, "*", ".*")
	s.handlers[path] = handler
}

// ZhinaoTestServer 创建一个模拟的360智脑 API 服务器
func (s *ServerTest) ZhinaoTestServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证 API Key
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + GetTestToken()

		if authHeader != expectedAuth {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":{"message":"Invalid API key","type":"invalid_request_error","code":"invalid_api_key"}}`))
			return
		}

		// 路由匹配
		for route, handler := range s.handlers {
			// 添加 ^ 和 $ 使路径匹配更精确
			pattern := regexp.MustCompile("^" + route + "$")
			if pattern.MatchString(r.URL.Path) {
				handler(w, r)
				return
			}
		}

		// 未找到匹配的路由
		http.Error(w, "Not Found", http.StatusNotFound)
	}))
}
