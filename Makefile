.PHONY: all build test bench clean examples lint coverage docs tools install-tools

# 默认目标
all: build test

# 构建
build:
	@echo "Building..."
	@go build -v ./...

# 构建命令行工具
tools:
	@echo "Building command line tools..."
	@go build -v ./cmd/gojson
	@go build -v ./cmd/jsonformat
	@go build -v ./cmd/jsonpath
	@go build -v ./cmd/jsonanalyze
	@go build -v ./cmd/jsonstream

# 安装命令行工具
install-tools:
	@echo "Installing command line tools..."
	@go install ./cmd/gojson
	@go install ./cmd/jsonformat
	@go install ./cmd/jsonpath
	@go install ./cmd/jsonanalyze
	@go install ./cmd/jsonstream

# 测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 基准测试
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./benchmarks

# 清理
clean:
	@echo "Cleaning..."
	@go clean
	@rm -f coverage.out
	@rm -f gojson jsonformat jsonpath jsonanalyze jsonstream

# 运行示例
examples:
	@echo "Running examples..."
	@for example in examples/*.go; do \
		echo "Running $$example"; \
		go run $$example; \
	done

# 代码检查
lint:
	@echo "Linting..."
	@golangci-lint run

# 代码覆盖率
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out


# 帮助
help:
	@echo "Available targets:"
	@echo "  all           - Build and test"
	@echo "  build         - Build the project"
	@echo "  tools         - Build command line tools"
	@echo "  install-tools - Install command line tools"
	@echo "  test          - Run tests"
	@echo "  bench         - Run benchmarks"
	@echo "  clean         - Clean build artifacts"
	@echo "  examples      - Run examples"
	@echo "  lint          - Run linter"
	@echo "  coverage      - Generate coverage report"
	@echo "  help          - Show this help message"
