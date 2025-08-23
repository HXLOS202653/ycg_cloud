# Makefile for Go project code quality and build automation
# 使用方法: make <target>

# 变量定义
GO_VERSION := 1.23
BINARY_NAME := ycg_cloud
MAIN_PATH := ./cmd/main.go
PKG := ./...
COVERAGE_FILE := coverage.out

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

## help: 显示帮助信息
.PHONY: help
help:
	@echo "$(BLUE)可用的命令:$(RESET)"
	@echo "$(GREEN)  make install-tools$(RESET)    - 安装所有代码质量工具"
	@echo "$(GREEN)  make format$(RESET)          - 格式化代码 (gofmt + goimports)"
	@echo "$(GREEN)  make lint$(RESET)            - 运行代码规范检查 (golangci-lint)"
	@echo "$(GREEN)  make security$(RESET)        - 运行安全扫描 (gosec)"
	@echo "$(GREEN)  make complexity$(RESET)      - 检查代码复杂度 (gocyclo)"
	@echo "$(GREEN)  make test$(RESET)            - 运行测试"
	@echo "$(GREEN)  make test-coverage$(RESET)   - 运行测试并生成覆盖率报告"
	@echo "$(GREEN)  make build$(RESET)           - 构建应用程序"
	@echo "$(GREEN)  make clean$(RESET)           - 清理构建文件"
	@echo "$(GREEN)  make check-all$(RESET)       - 运行所有代码质量检查"
	@echo "$(GREEN)  make pre-commit$(RESET)      - 运行pre-commit检查"
	@echo "$(GREEN)  make pre-push$(RESET)        - 运行pre-push检查"
	@echo "$(GREEN)  make deps$(RESET)            - 下载和整理依赖"
	@echo "$(GREEN)  make mod-tidy$(RESET)        - 整理go.mod文件"
	@echo "$(GREEN)  make run$(RESET)             - 运行应用程序"

## install-tools: 安装所有代码质量工具
.PHONY: install-tools
install-tools:
	@echo "$(BLUE)安装代码质量工具...$(RESET)"
	@echo "$(YELLOW)安装 golangci-lint...$(RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(YELLOW)安装 gosec...$(RESET)"
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "$(YELLOW)安装 gocyclo...$(RESET)"
	@go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@echo "$(YELLOW)安装 goimports...$(RESET)"
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(YELLOW)安装 gofumpt...$(RESET)"
	@go install mvdan.cc/gofumpt@latest
	@echo "$(GREEN)所有工具安装完成!$(RESET)"

## deps: 下载和整理依赖
.PHONY: deps
deps:
	@echo "$(BLUE)下载依赖...$(RESET)"
	@go mod download
	@go mod verify

## mod-tidy: 整理go.mod文件
.PHONY: mod-tidy
mod-tidy:
	@echo "$(BLUE)整理go.mod文件...$(RESET)"
	@go mod tidy

## format: 格式化代码
.PHONY: format
format:
	@echo "$(BLUE)格式化代码...$(RESET)"
	@echo "$(YELLOW)运行 gofmt...$(RESET)"
	@gofmt -s -w .
	@echo "$(YELLOW)运行 goimports...$(RESET)"
	@goimports -w -local ycg_cloud .
	@echo "$(YELLOW)运行 gofumpt...$(RESET)"
	@gofumpt -w .
	@echo "$(GREEN)代码格式化完成!$(RESET)"

## lint: 运行代码规范检查
.PHONY: lint
lint:
	@echo "$(BLUE)运行代码规范检查...$(RESET)"
	@golangci-lint run --config .golangci.yml
	@echo "$(GREEN)代码规范检查完成!$(RESET)"

## lint-fix: 运行代码规范检查并自动修复
.PHONY: lint-fix
lint-fix:
	@echo "$(BLUE)运行代码规范检查并自动修复...$(RESET)"
	@golangci-lint run --config .golangci.yml --fix
	@echo "$(GREEN)代码规范检查和修复完成!$(RESET)"

## security: 运行安全扫描
.PHONY: security
security:
	@echo "$(BLUE)运行安全扫描...$(RESET)"
	@gosec -fmt=json -out=gosec-report.json -stdout -verbose=text ./...
	@echo "$(GREEN)安全扫描完成! 报告已保存到 gosec-report.json$(RESET)"

## complexity: 检查代码复杂度
.PHONY: complexity
complexity:
	@echo "$(BLUE)检查代码复杂度...$(RESET)"
	@gocyclo -over 10 .
	@echo "$(GREEN)复杂度检查完成!$(RESET)"

## test: 运行测试
.PHONY: test
test:
	@echo "$(BLUE)运行测试...$(RESET)"
	@go test -v -race $(PKG)
	@echo "$(GREEN)测试完成!$(RESET)"

## test-coverage: 运行测试并生成覆盖率报告
.PHONY: test-coverage
test-coverage:
	@echo "$(BLUE)运行测试并生成覆盖率报告...$(RESET)"
	@go test -v -race -coverprofile=$(COVERAGE_FILE) $(PKG)
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "$(GREEN)测试覆盖率报告已生成: coverage.html$(RESET)"

## test-short: 运行短测试
.PHONY: test-short
test-short:
	@echo "$(BLUE)运行短测试...$(RESET)"
	@go test -short -v $(PKG)
	@echo "$(GREEN)短测试完成!$(RESET)"

## build: 构建应用程序
.PHONY: build
build:
	@echo "$(BLUE)构建应用程序...$(RESET)"
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)构建完成! 二进制文件: bin/$(BINARY_NAME)$(RESET)"

## build-linux: 为Linux构建应用程序
.PHONY: build-linux
build-linux:
	@echo "$(BLUE)为Linux构建应用程序...$(RESET)"
	@GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "$(GREEN)Linux构建完成! 二进制文件: bin/$(BINARY_NAME)-linux$(RESET)"

## build-windows: 为Windows构建应用程序
.PHONY: build-windows
build-windows:
	@echo "$(BLUE)为Windows构建应用程序...$(RESET)"
	@GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "$(GREEN)Windows构建完成! 二进制文件: bin/$(BINARY_NAME).exe$(RESET)"

## run: 运行应用程序
.PHONY: run
run:
	@echo "$(BLUE)运行应用程序...$(RESET)"
	@go run $(MAIN_PATH)

## clean: 清理构建文件
.PHONY: clean
clean:
	@echo "$(BLUE)清理构建文件...$(RESET)"
	@rm -rf bin/
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage.html
	@rm -f gosec-report.json
	@echo "$(GREEN)清理完成!$(RESET)"

## check-all: 运行所有代码质量检查
.PHONY: check-all
check-all: format lint security complexity test
	@echo "$(GREEN)所有代码质量检查完成!$(RESET)"

## pre-commit: 运行pre-commit检查
.PHONY: pre-commit
pre-commit:
	@echo "$(BLUE)运行pre-commit检查...$(RESET)"
	@$(MAKE) format
	@$(MAKE) lint
	@$(MAKE) test-short
	@echo "$(GREEN)Pre-commit检查完成!$(RESET)"

## pre-push: 运行pre-push检查
.PHONY: pre-push
pre-push:
	@echo "$(BLUE)运行pre-push检查...$(RESET)"
	@$(MAKE) check-all
	@$(MAKE) build
	@echo "$(GREEN)Pre-push检查完成!$(RESET)"

## docker-build: 构建Docker镜像
.PHONY: docker-build
docker-build:
	@echo "$(BLUE)构建Docker镜像...$(RESET)"
	@docker build -t $(BINARY_NAME):latest .
	@echo "$(GREEN)Docker镜像构建完成!$(RESET)"

## docker-run: 运行Docker容器
.PHONY: docker-run
docker-run:
	@echo "$(BLUE)运行Docker容器...$(RESET)"
	@docker run --rm -p 8080:8080 $(BINARY_NAME):latest

## version: 显示Go版本信息
.PHONY: version
version:
	@echo "$(BLUE)Go版本信息:$(RESET)"
	@go version
	@echo "$(BLUE)项目模块:$(RESET)"
	@go list -m

## env: 显示Go环境信息
.PHONY: env
env:
	@echo "$(BLUE)Go环境信息:$(RESET)"
	@go env

# 检查工具是否安装
check-tool = $(shell command -v $(1) 2> /dev/null)

## check-tools: 检查所需工具是否已安装
.PHONY: check-tools
check-tools:
	@echo "$(BLUE)检查工具安装状态...$(RESET)"
	@echo -n "golangci-lint: "
	@if [ "$(call check-tool,golangci-lint)" ]; then echo "$(GREEN)✓ 已安装$(RESET)"; else echo "$(RED)✗ 未安装$(RESET)"; fi
	@echo -n "gosec: "
	@if [ "$(call check-tool,gosec)" ]; then echo "$(GREEN)✓ 已安装$(RESET)"; else echo "$(RED)✗ 未安装$(RESET)"; fi
	@echo -n "gocyclo: "
	@if [ "$(call check-tool,gocyclo)" ]; then echo "$(GREEN)✓ 已安装$(RESET)"; else echo "$(RED)✗ 未安装$(RESET)"; fi
	@echo -n "goimports: "
	@if [ "$(call check-tool,goimports)" ]; then echo "$(GREEN)✓ 已安装$(RESET)"; else echo "$(RED)✗ 未安装$(RESET)"; fi
	@echo -n "gofumpt: "
	@if [ "$(call check-tool,gofumpt)" ]; then echo "$(GREEN)✓ 已安装$(RESET)"; else echo "$(RED)✗ 未安装$(RESET)"; fi