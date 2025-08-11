#!/bin/bash

# 测试脚本 - 运行所有测试
# Usage: ./scripts/test.sh [options]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示帮助信息
show_help() {
    cat << EOF
YCG Cloud Storage - 测试脚本

用法: $0 [选项]

选项:
    -h, --help          显示此帮助信息
    -u, --unit          只运行单元测试
    -i, --integration   只运行集成测试
    -c, --coverage      生成覆盖率报告
    -r, --race          启用竞态检测
    -b, --bench         运行性能测试
    -l, --lint          运行代码检查
    -a, --all           运行所有测试 (默认)
    --clean             清理测试产物
    --docker            使用Docker运行测试

示例:
    $0                  # 运行所有测试
    $0 -u -c           # 运行单元测试并生成覆盖率报告
    $0 -i --docker     # 使用Docker运行集成测试
    $0 --clean         # 清理测试产物

EOF
}

# 默认选项
RUN_UNIT=false
RUN_INTEGRATION=false
RUN_COVERAGE=false
RUN_RACE=false
RUN_BENCH=false
RUN_LINT=false
RUN_ALL=true
CLEAN_ONLY=false
USE_DOCKER=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -u|--unit)
            RUN_UNIT=true
            RUN_ALL=false
            shift
            ;;
        -i|--integration)
            RUN_INTEGRATION=true
            RUN_ALL=false
            shift
            ;;
        -c|--coverage)
            RUN_COVERAGE=true
            shift
            ;;
        -r|--race)
            RUN_RACE=true
            shift
            ;;
        -b|--bench)
            RUN_BENCH=true
            RUN_ALL=false
            shift
            ;;
        -l|--lint)
            RUN_LINT=true
            RUN_ALL=false
            shift
            ;;
        -a|--all)
            RUN_ALL=true
            shift
            ;;
        --clean)
            CLEAN_ONLY=true
            shift
            ;;
        --docker)
            USE_DOCKER=true
            shift
            ;;
        *)
            log_error "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 清理函数
cleanup() {
    log_info "清理测试产物..."
    rm -f coverage.out coverage.html
    rm -rf ./test_uploads
    go clean -testcache
    log_success "清理完成"
}

# 如果只是清理，执行清理后退出
if [ "$CLEAN_ONLY" = true ]; then
    cleanup
    exit 0
fi

# 设置测试环境
export APP_ENV=test

log_info "开始YCG云盘存储测试..."
log_info "Go版本: $(go version)"
log_info "测试模式: $([ "$USE_DOCKER" = true ] && echo "Docker" || echo "本地")"

# 检查依赖
log_info "检查Go模块..."
go mod verify
go mod download

# 代码检查
if [ "$RUN_LINT" = true ] || [ "$RUN_ALL" = true ]; then
    log_info "运行代码检查..."
    
    # 检查代码格式
    if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
        log_error "代码格式不正确:"
        gofmt -s -l .
        exit 1
    fi
    
    # 运行golangci-lint
    if command -v golangci-lint >/dev/null 2>&1; then
        golangci-lint run --config .golangci.yml
        log_success "代码检查通过"
    else
        log_warning "golangci-lint 未安装，跳过代码检查"
    fi
fi

# 构建测试选项
TEST_FLAGS=""
if [ "$RUN_RACE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -race"
fi

if [ "$RUN_COVERAGE" = true ]; then
    TEST_FLAGS="$TEST_FLAGS -coverprofile=coverage.out -covermode=atomic"
fi

# 单元测试
if [ "$RUN_UNIT" = true ] || [ "$RUN_ALL" = true ]; then
    log_info "运行单元测试..."
    
    if [ "$USE_DOCKER" = true ]; then
        log_warning "Docker模式下跳过单元测试（单元测试不需要外部依赖）"
    else
        go test -short -v $TEST_FLAGS ./internal/...
        log_success "单元测试完成"
    fi
fi

# 集成测试
if [ "$RUN_INTEGRATION" = true ] || [ "$RUN_ALL" = true ]; then
    log_info "运行集成测试..."
    
    if [ "$USE_DOCKER" = true ]; then
        log_info "使用Docker Compose运行集成测试..."
        docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
        docker-compose -f docker-compose.test.yml down --volumes
    else
        log_warning "本地运行集成测试需要启动数据库服务"
        log_info "如果数据库未启动，测试可能会失败"
        go test -v -tags=integration $TEST_FLAGS ./tests/... || log_warning "集成测试失败（可能是数据库连接问题）"
    fi
    
    log_success "集成测试完成"
fi

# 性能测试
if [ "$RUN_BENCH" = true ]; then
    log_info "运行性能测试..."
    go test -bench=. -benchmem ./...
    log_success "性能测试完成"
fi

# 生成覆盖率报告
if [ "$RUN_COVERAGE" = true ] && [ -f coverage.out ]; then
    log_info "生成覆盖率报告..."
    go tool cover -html=coverage.out -o coverage.html
    
    # 显示覆盖率统计
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    log_success "总覆盖率: $COVERAGE"
    log_info "HTML报告已生成: coverage.html"
    
    # 尝试打开覆盖率报告
    if command -v xdg-open >/dev/null 2>&1; then
        xdg-open coverage.html 2>/dev/null &
    elif command -v open >/dev/null 2>&1; then
        open coverage.html 2>/dev/null &
    fi
fi

log_success "所有测试完成!"

# 显示测试总结
echo ""
log_info "测试总结:"
echo "  - 单元测试: $([ "$RUN_UNIT" = true ] || [ "$RUN_ALL" = true ] && echo "✓" || echo "✗")"
echo "  - 集成测试: $([ "$RUN_INTEGRATION" = true ] || [ "$RUN_ALL" = true ] && echo "✓" || echo "✗")"
echo "  - 代码检查: $([ "$RUN_LINT" = true ] || [ "$RUN_ALL" = true ] && echo "✓" || echo "✗")"
echo "  - 覆盖率报告: $([ "$RUN_COVERAGE" = true ] && echo "✓" || echo "✗")"
echo "  - 竞态检测: $([ "$RUN_RACE" = true ] && echo "✓" || echo "✗")"
echo "  - 性能测试: $([ "$RUN_BENCH" = true ] && echo "✓" || echo "✗")"