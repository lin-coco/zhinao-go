# 360智脑 Go SDK Makefile

.PHONY: test test-verbose test-coverage test-unit test-mock build clean lint help

# 默认目标
.DEFAULT_GOAL := help

# 运行测试
test:
	@echo "Running tests..."
	@go test -v $(shell go list ./... | grep -v /examples)

# 运行测试（详细输出）
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out $(shell go list ./... | grep -v /examples)
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 仅运行单元测试（不包括 examples）
test-unit:
	@echo "Running unit tests..."
	@go test -v .

# 仅运行 Mock 测试
test-mock:
	@echo "Running mock tests..."
	@go test -v -run Mock .

# 构建项目
build:
	@echo "Building project..."
	@go build -v ./...

# 清理生成的文件
clean:
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@go clean

# 代码检查
lint:
	@echo "Running linters..."
	@go vet ./...
	@gofmt -l .

# 格式化代码
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 检查依赖
deps:
	@echo "Checking dependencies..."
	@go mod tidy
	@go mod verify

# 显示帮助信息
help:
	@echo "360智脑 Go SDK - Makefile 命令"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用目标:"
	@echo "  test           - 运行所有测试"
	@echo "  test-verbose   - 运行测试（详细输出）"
	@echo "  test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  test-unit      - 仅运行单元测试"
	@echo "  test-mock      - 仅运行 Mock 测试"
	@echo "  build          - 构建项目"
	@echo "  clean          - 清理生成的文件"
	@echo "  lint           - 运行代码检查"
	@echo "  fmt            - 格式化代码"
	@echo "  deps           - 检查和整理依赖"
	@echo "  help           - 显示此帮助信息"
